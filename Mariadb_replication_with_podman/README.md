### MariaDB replication in Podman

This task was performed a declarative way and in the Podman's rootless containers only.

Image official docker.io/library/mariadb:11.7-ubi

Get Podman version

```
podman version
Client:       Podman Engine
Version:      4.6.2
API Version:  4.6.2
Go Version:   go1.19.12
Built:        Mon Aug 28 22:38:31 2023
OS/Arch:      linux/amd64
```

#### Configure primary server

Create master-init.sql file to initialize the user for replication and a new database.

```
CREATE USER 'repluser'@'%' IDENTIFIED BY 'secret';
GRANT REPLICATION SLAVE ON *.* TO 'repluser'@'%';
CREATE DATABASE primary_db;

```

Create config file master-conf.cnf to enable binary logging, set server_id to the master
 and set port=3307 to avoid conflict equal ports with the slave server.

```
[mariadb]
log-bin                        
server_id=7000                

[mysqld]
port=3307

```

To connect .sql file to podman create podman volume and copy this file to it

`localhost$ podman volume create vol-master-init`

`localhost$ cp master-init.sql ~/.local/share/containers/storage/volumes/vol-master-init/_data/`

The same procedure to the .cnf file

`localhost$ podman volume create vol-master-conf`

`localhost$ cp master-conf.cnf ~/.local/share/containers/volumes/vol-master-conf/_data/`

Create mariadb-secret.yaml file to store sensitive data - password of the root user

```
apiVersion: v1
data:
  password: c2VjcmV0 # get this string in vim ':r!printf "secret" | base64'
kind: Secret
metadata:
  creationTimestamp: null
  name: master-password

```

And create podman secret based on this file.

`localhost$ podman kube play mariadb-secret.yaml`

#### Configure secondary server

Create slave-init.sql file

```

CHANGE MASTER TO
   MASTER_HOST='128.0.0.1', # both containers will be in one pod and share the localhost 
   MASTER_USER='repluser',  #network and have equal ip addresses
   MASTER_PASSWORD='secret',
   MASTER_PORT=3307,        # 3307 from prime server config
   MASTER_CONNECT_RETRY=10;

```

Create podman volume to bind this content to the slave container.

`localhost$ podman volume create vol-slave-init`

Copy the init file to a volume.

`localhost$ cp slave-init.sql ~/.local/share/containers/storage/volumes/vol-slave-init/_data/`

Create slave-conf.cnf file to enable binary logging and set server_id not equal to master's server_id

```
[mariadb]
log-bin                         
server_id=4999

```
Create podman volume and copy .cnf file to it.

`localhost$ podman volume create vol-slave-conf`

`localhost$ cp slave-conf.cnf ~/.local/share/containers/storage/volumes/vol-slave-conf/_data/`

After each of these operations it possible to ofcourse checking results with:

__podman volume <secret> ls__ or __podman volume <secret> inspect vol-name <secret-name>__


#### Create Kubernetes manifest yaml file to put it all together  


```

apiVersion: v1
kind: Pod
metadata:
  name: mariadb # the name of Podman pod
spec: 
  volumes:
    - name: vol-master-init  #the list of the volumes were created before
      emptyDir:              #aim to set data to containers
        sizeLimit: 100Mi     
    - name: vol-master-conf
      emptyDir:
        sizeLimit: 100Mi
    - name: vol-slave-conf
      emptyDir:
        sizeLimit: 50Mi
    - name: vol-slave-init
      emptyDir:
        sizeLimit: 50Mi
   containers:
    - name: master           #primary container full name will be mariadb-master
      image: docker.io/library/mariadb:11.7-ubi
      env:
        - name: MARIADB_ROOT_PASSWORD #environment MariaDB
          valueFrom:
            secretKeyRef:
              name: master-password  #the name of podman secret 
              key: password
      volumeMounts:
        - mountPath: "/etc/mysql/conf.d" #config directory in the master container
          name: vol-master-conf
        - mountPath: "/docker-entrypoint-initdb.d" #entry point master container
          name: vol-master-init
    - name: slave #the full name of the secondary container will be mariadb-slave
      image: docker.io/library/mariadb:11.7-ubi
      env:
        - name: MARIADB_ALLOW_EMPTY_ROOT_PASSWORD #also possible to set podman secret here
          value: True 
      volumeMounts:
        - mountPath: "/docker-entrypoint-initdb.d"
          name: vol-slave-init
        - mountPath: "/etc/mysql/conf.d"
          name: vol-slave-conf

```

#### Run creating pod

`localhost$ podman kube play mariadb-pod.yaml`

Check the result

```
localhost$ podman ps
CONTAINER ID  IMAGE                                    COMMAND     CREATED         STATUS        PORTS       NAMES
0fd54239d78a  localhost/podman-pause:4.6.2-1693251511              10 seconds ago  Up 5 seconds              4b7559660459-infra
eea2485ea9e4  docker.io/library/mariadb:11.7-ubi       mariadbd    8 seconds ago   Up 4 seconds              mariadb-master
2aee44090e30  docker.io/library/mariadb:11.7-ubi       mariadbd    6 seconds ago   Up 5 seconds              mariadb-slave

```

get logs for both containers

`localhost$ podman logs mariadb-master `

```
............................
2025-02-22 11:24:43 0 [Note] mariadbd: ready for connections.
Version: '11.7.2-MariaDB-log'  socket: '/run/mariadb/mariadb.sock'  port: 3307  MariaDB Server
2025-02-22 11:24:45 4 [Note] Start binlog_dump to slave_server(1), pos(, 4), using_gtid(1), gtid('')

```

`localhost$ podman logs mariadb-slave`

```
..........................

2025-02-22 11:24:45 0 [Note] mariadbd: ready for connections.
Version: '11.7.2-MariaDB'  socket: '/run/mariadb/mariadb.sock'  port: 3306  MariaDB Server
2025-02-22 11:24:45 5 [Note] Slave SQL thread initialized, starting replication in log 'FIRST' at position 4, relay log './mariadb-relay-bin.000001' position: 4; GTID position ''
2025-02-22 11:24:45 4 [Note] Slave I/O thread: connected to master 'repluser@127.0.0.1:3307',replication starts at GTID position ''
2025-02-22 11:24:42+00:00 [Note] [Entrypoint]: Stopping temporary server
2025-02-22 11:24:42+00:00 [Note] [Entrypoint]: Temporary server stopped

```

check the binary loging and replication status on the primary container

```
podman exec -it mariadb-master mariadb -uroot -psecret -e 'show binary logs;'
+-----------------------+-----------+
| Log_name              | File_size |
+-----------------------+-----------+
| my-mariadb-bin.000001 |       820 |
| my-mariadb-bin.000002 |       347 |
+-----------------------+-----------+

```

```

podman exec -it mariadb-master mariadb -uroot -psecret -e 'show master status;'
+-----------------------+----------+--------------+------------------+
| File                  | Position | Binlog_Do_DB | Binlog_Ignore_DB |
+-----------------------+----------+--------------+------------------+
| my-mariadb-bin.000002 |      347 |              |                  |
+-----------------------+----------+--------------+------------------+

```

check the appearing new database and replication status on thesecondary container

```

podman exec -it mariadb-slave mariadb -uroot -e 'show databases;'
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| primary_db         | <-- new database is here now
| sys                |
+--------------------+

```

```

podman exec -it mariadb-slave mariadb -uroot -e 'show replica status\G'
*************************** 1. row ***************************
                Slave_IO_State: Waiting for master to send event
                   Master_Host: 127.0.0.1
                   Master_User: repluser
                   Master_Port: 3307
                 Connect_Retry: 10
               Master_Log_File: my-mariadb-bin.000002
           Read_Master_Log_Pos: 347
                Relay_Log_File: mariadb-relay-bin.000003
                 Relay_Log_Pos: 651
         Relay_Master_Log_File: my-mariadb-bin.000002
              Slave_IO_Running: Yes
             Slave_SQL_Running: Yes
               Replicate_Do_DB: 
           Replicate_Ignore_DB: 
            Replicate_Do_Table: 
        Replicate_Ignore_Table: 
       Replicate_Wild_Do_Table: 
   Replicate_Wild_Ignore_Table: 
                    Last_Errno: 0
                    Last_Error: 
                  Skip_Counter: 0
           Exec_Master_Log_Pos: 347
               Relay_Log_Space: 1807
               Until_Condition: None
                Until_Log_File: 
                 Until_Log_Pos: 0
            Master_SSL_Allowed: Yes
            Master_SSL_CA_File: 
            Master_SSL_CA_Path: 
               Master_SSL_Cert: 
             Master_SSL_Cipher: 
                Master_SSL_Key: 
         Seconds_Behind_Master: 0
 Master_SSL_Verify_Server_Cert: Yes
                 Last_IO_Errno: 0
                 Last_IO_Error: 
                Last_SQL_Errno: 0
                Last_SQL_Error: 
   Replicate_Ignore_Server_Ids: 
              Master_Server_Id: 5000
                Master_SSL_Crl: 
            Master_SSL_Crlpath: 
                    Using_Gtid: Slave_Pos
                   Gtid_IO_Pos: 0-3000-3
       Replicate_Do_Domain_Ids: 
   Replicate_Ignore_Domain_Ids: 
                 Parallel_Mode: optimistic
                     SQL_Delay: 0
           SQL_Remaining_Delay: NULL
       Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
              Slave_DDL_Groups: 3
Slave_Non_Transactional_Groups: 0
    Slave_Transactional_Groups: 0
          Replicate_Rewrite_DB: 

```

finally stop pod and containers

```

localhost$ podman kube play mariadb-pod.yaml --down
Pods stopped:
4b75596604590d2edaaf3fd1e4b5ed2bdd89736480d3c9c6959ac491624d156e
Pods removed:
4b75596604590d2edaaf3fd1e4b5ed2bdd89736480d3c9c6959ac491624d156e
Secrets removed:
Volumes removed:

```