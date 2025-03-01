### PostgreSQL logical replication with Podman

There will be two rootless containers that implement Docker PostgreSQL official image.

Podman version

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

Create podman volume to store the init script

`localhost$ podman volume create logical-master`

Bash script performs the next tasks

add a line at the end of pg_hba.conf

`echo "host replication replicator 127.0.0.1/32 md5" >> /var/lib/postgresql/data/pg_hba.conf`

set some settings to replication configure, ports in primary and secondary must not be equal 

```
    create user replicator with replication password 'secret';
    alter system set listen_addresses to '*';
    alter system set wal_level to logical;
    alter system set port to 2345;
```

create a database to replicate and the publication all tables of it

```
    create database sourcedb;
    \c sourcedb
    create table tb(id integer not null primary key, name varchar(64));
    grant select on all tables in schema public to replicator;
    create publication publication_one for all tables;

```

copy this script to the logical-master volume

`localhost$ cp logical-master ~/.local/share/containers/storage/volumes/logical-master/_data/`


#### Configure secondary server

create podman volume to store init slave script

`localhost$ podman volume create init-slave`

Bash script to perform the next tasks

To make sure running PostgreSQL on the secondary server after PostgreSQL on the primary one

`sleep 5s` may be less or more seconds

create a database and table with some schema as on the primary server

```
    create database destinationdb;
    \c destinationdb
    create table tb(id integer not null primary key, name varchar(64));
```

create the subscription to the publication sourcedb database on the primary server

```
    create subscription slavesub connection 'user=replicator password=secret host=127.0.0.1 port=2345 dbname=sourcedb' \
    publication publication_one;
```

copy this script to the volume init-slave

`localhost$ cp init-slave.sh ~/.local/share/containers/storage/volumes/init-slave/_data/`

#### Create a secret for the sensitive data

secret.yaml file for both primary and secondary containers

```
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  name: logical-replication
data:
  password: c2VjcmV0

```
create podman secret

```
podman kube play secret.yaml 
Secrets:
1a7a22b0f8bb8e5303f7d4928

```

#### Put it all together with Kubernetes manifest postgres-pod.yaml file

podman will create pod with two containers postgres-master and postgres-slave inside;

According to [image documentation](https://hub.docker.com/_/postgres) these scripts have to be mounted to the /docker-entrypoint-initdb.d directory

Run containers

```
localhost$ podman kube play postgres-pod.yaml
Pod:
cef0916a6cd8ecbe6dc37a8c7261b8f9c933e51e6c2e8028aa135a2bb2274e95
Containers:
98e06105fb06d3cca2f3fa37bceb4492145709db3473539a7eaecebc33987559
b8e6a962578da8016807cce9f5ce6a2d79b1cab5c7322c1b709259b58335dac4
```
check the state of containers

```
localhost$ podman ps -a
CONTAINER ID  IMAGE                                    COMMAND     CREATED         STATUS        PORTS       NAMES
889bec3dd808  localhost/podman-pause:4.6.2-1693251511              16 seconds ago  Up 7 seconds              cef0916a6cd8-infra
98e06105fb06  docker.io/library/postgres:17-bookworm   postgres    13 seconds ago  Up 6 seconds              postgres-master
b8e6a962578d  docker.io/library/postgres:17-bookworm   postgres    10 seconds ago  Up 5 seconds              postgres-slave

```

check the logs of the containers

`localhost$ podman logs postgres-master`

`localhost$ podman logs postgres-slave`

show that logical replication has been done successfully on the source and destination databases

insert some content into the sourcedb database

`localhost$ podman exec -it postgres-master psql -U postgres -p 2345`

```
postgres=# \c sourcedb
You are now connected to database "sourcedb" as user "postgres".

```
```
sourcedb=# \dt
        List of relations
 Schema | Name | Type  |  Owner   
--------+------+-------+----------
 public | tb   | table | postgres
(1 row)

sourcedb=# insert into tb values(1,'first'),(2,'second');
INSERT 0 2
sourcedb=# select * from tb;
 id |  name  
----+--------
  1 | first
  2 | second
(2 rows)
```

check replication result

`podman exec -it postgres-slave -U postgres`

```
postgres=# \c destinationdb 
You are now connected to database "destinationdb" as user "postgres".
destinationdb=# \dt
        List of relations
 Schema | Name | Type  |  Owner   
--------+------+-------+----------
 public | tb   | table | postgres
(1 row)

destinationdb=# select * from tb;
 id |  name  
----+--------
  1 | first
  2 | second
(2 rows)
```

stop and remove containers

```
localhost$ podman kube play postgres-pod.yaml --down
Pods stopped:
63bb40f8803f3f867d1e397eb79802b134209b89248f34b38f20c735aa03acbc
Pods removed:
63bb40f8803f3f867d1e397eb79802b134209b89248f34b38f20c735aa03acbc
Secrets removed:
Volumes removed:

```
