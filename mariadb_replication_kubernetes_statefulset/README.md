## MariaDB replication with Kubernetes StatefulSet

Operating System: Fedora Linux 37 (Workstation Edition)

Minikube version

```
localhost$ minikube version
minikube version: v1.33.0
```

create a namespace and make it default to work with

`localhost$ kubectl create namespace mariadb-ns`

`localhost$ kubectl config set-context --current --namespace mariadb-ns`

### Create ConfigMap

ConfigMap consists of four configure and executing files. Two of them configure the primary server

```
  master.cnf: |
    [mariadb]
    log-bin
    log-basename=my-mariadb
  master.sql: |
    CREATE USER 'replicator'@'%' IDENTIFIED BY 'secret';
    GRANT REPLICATION SLAVE ON *.* TO 'replicator'@'%';
    CREATE DATABASE replicatest;

```
Configure secondary server(s)

```
  slave.cnf: |
    [mariadb]
    relay-log=mysqld-relay-bin
    log-basename=my-mariadb
  slave.sql: |
      CHANGE MASTER TO
      #$(statefulset-number).$(headlesservice).$(namespace).$(domain-cluster)  to correct work DNS
      MASTER_HOST='mariadb-statefulset-0.mariadb-service.mariadb-ns.svc.cluster.local',
      MASTER_USER='replicator',
      MASTER_PASSWORD='secret',
      MASTER_USE_GTID=slave_pos,
      MASTER_CONNECT_RETRY=10;

```

### Create Headless Service

To connect directly to a Pod without proxy and load balancing 

```
apiVersion: v1
kind: Service
metadata:
  name: mariadb-service
  labels:
    app: mariadb
spec:
  ports:
    - port: 3306
      name: mariadb-port
  clusterIP: None #pay atantion
  selector:
    app: mariadb
```
### Create Secret with Kustomize

Create kustomization.yaml file with content

```
secretGenerator:
  - name: mariadb-creds
    literals:
      - password=secret
generatorOptions:
  disableNameSuffixHash: True
  labels:
    type: generated
  annotations: 
    note: secret-generated

```

### Create Statefulset

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mariadb-statefulset
spec:
  serviceName: "mariadb-service" #headless service
  selector:
    matchLabels: 
      app: mariadb 
  replicas: 3
  template:
    metadata:
      labels:
        app: mariadb
    spec:
      initContainers: #init container to start MariaDB on the servers
        - name: init-mariadb
          image: mariadb
          volumeMounts:
            - name: mariadb-global-config #store files from ConfigMap
              mountPath: /mnt/conf.d      #created by kubelet
            - name: mariadb-cnf
              mountPath: /etc/mysql/conf.d
            - name: mariadb-init
              mountPath: /docker-entrypoint-initdb.d
          command:
           - bash
           - "-c"
           - |
             set -ex
             [[ $(hostname) =~ -([0-9]+)$ ]] || exit 1
             suffix=${BASH_REMATCH[1]} #return 0 to the primary server and 1,2...n to replica servers
             #may be 1000 or any other number, containerId range: 1 to 4294967295 
             containerId=$((1000+$suffix))
             if [[ $suffix -eq 0 ]] 
             then #primary server configuration
             cp /mnt/conf.d/master.cnf /etc/mysql/conf.d/server.cnf
             cp /mnt/conf.d/master.sql /docker-entrypoint-initdb.d/
             else  #secondary server(s) configuration
             cp /mnt/conf.d/slave.cnf /etc/mysql/conf.d/server.cnf
             cp /mnt/conf.d/slave.sql /docker-entrypoint-initdb.d/
             fi
             #add this line to both primary and secondary servers' config
             echo "server-id=$containerId" >> /etc/mysql/conf.d/server.cnf
      restartPolicy: Always
      containers: #working containers    
        - name: mariadb
          image: mariadb
          env:
            - name: MARIADB_ROOT_PASSWORD #required environment
              valueFrom:
                secretKeyRef:
                  name: mariadb-creds
                  key: password
          ports:
            - containerPort: 3306
              name: mariadb-port
          volumeMounts:
            # replication works without datadir
            - name: datadir #persistent storage to save data when restarting the database
              mountPath: /var/lib/mysql/
            - name: mariadb-cnf #volume implements when starting containers to configure
              mountPath: /etc/mysql/conf.d
            - name: mariadb-init
              mountPath: /docker-entrypoint-initdb.d #volume implements when starting containers to init
      volumes:
        - name: mariadb-global-config
          configMap:
            name: mariadb-config
            defaultMode: 0777 #.sql files executable .cnf writabel and readable
        # ConfigMap does not allow to execute and write. Copy config files to emptyDir volumes
        - name: mariadb-cnf
          emptyDir: {}
        - name: mariadb-init
          emptyDir: {}
  volumeClaimTemplates:
  - metadata: 
      name: datadir
    spec:
      accessModes: 
      - ReadWriteOnce
      resources:
        requests:
          storage: 300M

```

### Put it all together with Kustomize

Add these lines to the kustomization.yaml

```
...
resources:
  - statefulset.yaml
  - configmap.yaml
  - service.yaml

```
Store all files in one directory and run

```
localhost$ kubectl apply -k .
configmap/mariadb-config created
secret/mariadb-creds created
service/mariadb-service created
statefulset.apps/mariadb-statefulset created

```

Check the result with a set of commands

```
localhost$ kubectl logs mariadb-statefulset-(number)
localhost$ kubectl exec -it mariadb-statefulset-(number) -- mariadb -uroot -psecret -e 'show master(slave) status\G'
localhost$ kubectl exec -it mariadb-statefulset-(number) -- mariadb -uroot -psecret -e 'show variables like "server_id";'
localhost$ kubectl exec -it mariadb-statefulset-(number) -- mariadb -uroot -psecret -e 'show variables like "log_bin";'
localhost$ kubectl exec -it mariadb-statefulset-(number) -- mariadb -uroot -psecret -e 'show databases;'
```


