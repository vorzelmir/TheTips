### PostgreSQL streaming replication

On the primary and secondary servers must be equal PostgreSQL version.

On the both servers, it has to enable firewalld permissions

`primary<secondary># firewall-cmd --add-port=5432/tcp`

#### Configure primary server

primary server OS and PostgreSQL version

```
primary$ hostnamectl | grep -i 'operating system' | cut -d ':' -f 2
 Red Hat Enterprise Linux 9.4 (Plow)
```
```
primary$ psql -U postgres -c 'select version();'
                                                 version                                                  
----------------------------------------------------------------------------------------------------------
 PostgreSQL 17.4 on x86_64-pc-linux-gnu, compiled by gcc (GCC) 11.5.0 20240719 (Red Hat 11.5.0-5), 64-bit
(1 row)
```

create user(role) with the ability login and replication setting

`primary$ psql -U postgres -c 'create role replicator login password 'secret';`

`primary$ psql -U postgres -c 'alter role replicator replication;'`

`primary$ psql -U postgres -c '\du'`

create config file primary-replica.conf in the $PGDATA

```
wal_level=replica
max_wal_senders=10
wal_keep_size=100
listen_addresses='*'
hot_standby=on
archive_mode=on
archive_command='/usr/bin/true'

```

and add a reference to this file inside postgresql.conf with directive  __include__

```
......
include='primary-replica.conf'
......
```

add this string to the pg_hba.conf, where replica-server.ip is a DNS from /etc/hosts

`host   replication      replicator     replica-server.ip    md5`


after that, it needs to restart PostgreSQL server

`primary# systemctl restart postgresql-17`

now it is possible to get some setting variables, for instance

```
primary$ psql -U postgres -c 'show wal_level;'
 wal_level 
-----------
 replica
(1 row)

```

#### Configure secondary server

get OS system and PostgreSQL version

```
secondary$ hostnamectl | grep -i 'operating system'
  Operating System: Oracle Linux Server 9.5

```

```
psql -U postgres -c 'select version();'
                                                 version                                                  
----------------------------------------------------------------------------------------------------------
 PostgreSQL 17.4 on x86_64-pc-linux-gnu, compiled by gcc (GCC) 11.5.0 20240719 (Red Hat 11.5.0-5), 64-bit
(1 row)

```

Remove $PGDATA directory, in this case it is /var/lib/pgsql/17/data

`secondary# rm -rf /var/lib/pgsql/17/data`

Make a backup primary server, where primary-server.ip is a DNS from /etc/hosts

`secondary$ sudo -u postgres pg_basebackup -h primary-server.ip -U replicator -p 5432 -D /var/lib/pgsql/17/data -P -Xs -R`

Start the secondary PostgreSQL server

`secondary# systemctl start postgresql-17`

#### Check results

On the primary server create database

`primary$ psql -U postgres -c 'creaete database test_replica;'`

This database is now on the secondary server

`secondary$ psql -U postgres -c '\l'`

```
                                                       List of databases
     Name     |  Owner   | Encoding | Locale Provider |   Collate   |    Ctype    | Locale | ICU Rules |   Access privileges   
--------------+----------+----------+-----------------+-------------+-------------+--------+-----------+-----------------------
 postgres     | postgres | UTF8     | libc            | en_US.UTF-8 | en_US.UTF-8 |        |           | 
 template0    | postgres | UTF8     | libc            | en_US.UTF-8 | en_US.UTF-8 |        |           | =c/postgres          +
              |          |          |                 |             |             |        |           | postgres=CTc/postgres
 template1    | postgres | UTF8     | libc            | en_US.UTF-8 | en_US.UTF-8 |        |           | =c/postgres          +
              |          |          |                 |             |             |        |           | postgres=CTc/postgres
 test_replica | postgres | UTF8     | libc            | en_US.UTF-8 | en_US.UTF-8 |        |           | 
(5 rows)

```
