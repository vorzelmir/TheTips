### backup postgres database

got the dump .sql file [from here](https://neon.tech/docs/import/import-sample-data#world-happiness-index)

create database and connect to it

`postgres=# create database world_happiness;`

`postgres=# \c world_happiness`

load data from dump .sql file in this database

`world_happiness=# \i happiness_index.sql`

get list of tables

`world_happiness=# \d`

```
                 List of relations
 Schema |         Name          |   Type   | Owner  
--------+-----------------------+----------+--------
 public | 2019                  | table    | tomcat
 public | 2019_overall_rank_seq | sequence | tomcat
(2 rows)
```

create role to backup

`world_happiness=# \c postgres`

`postgres=# create role backup`

`postgres=# grant connect on database world_happiness to backup`

add some permissions to read from and write to database

`postgres=# alter role backup in database world_happiness set role pg_read_all_data;`

`postgres=# alter role backup in database world_happiness set role pg_write_all_data;`


get the result

`postgres=# \l`

```
...
 world_happiness | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |      | libc        | =Tc/tomcat       +
                 |        |          |             |             |      |             | tomcat=CTc/tomcat+
                 |        |          |             |             |      |             | backup=c/tomcat
...

```

#### logical backup

with **pg_dump**

`localhost$ pg_dump -U backup -d world_happiness -f world.dump.sql`

and use Insert instead of default Copy

`localhost$ pg_dump --insert -U backup -d world_happiness -f world.dump.sql`

to restore dump databases need to create new databases 

`localhost$ psql -Utomcat -c 'create database world_from_dump;`

then restore to this database

`localhost pg_dump --insert -f world.dump.sql -d world_from_dump`

or from psql

`world_from_dump=>\i world.dump.sql`

get the existing databases

`localhost$ psql -Utomcat -l`

```
                                                    List of databases
         Name      | Owner  | Encoding |   Collate   |    Ctype    | ICU Locale | Locale Provider | Access privileges 
-------------------+--------+----------+-------------+-------------+------------+-----------------+-------------------
 periodic_table    | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | =Tc/tomcat       +
                   |        |          |             |             |            |                 | tomcat=CTc/tomcat+
                   |        |          |             |             |            |                 | backup=c/tomcat
 postgres          | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | 
 template0         | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | =c/tomcat        +
                   |        |          |             |             |            |                 | tomcat=CTc/tomcat
 template1         | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | =c/tomcat        +
                   |        |          |             |             |            |                 | tomcat=CTc/tomcat
 world_from_dump   | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | 
 world_happiness   | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | =Tc/tomcat       +
                   |        |          |             |             |            |                 | tomcat=CTc/tomcat+
                   |        |          |             |             |            |                 | backup=c/tomcat
(6 rows)

```
another way to backup to customer binary format

`localhost$ pg_dump -U backup -Fc --create -f periodic.backup periodic_table`

drop database

`localhost$ psql -Utomcat -c 'drop database periodic_table;'`

restore from this format with **pg_restore**

`localhost$ pg_restore -C -d postgres periodic.backup`

check the result 

`localhost$ psql -d periodic_table -c '\dt'`

```
            List of relations
 Schema |      Name      | Type  | Owner  
--------+----------------+-------+--------
 public | periodic_table | table | tomcat
(1 row)
```


#### physical backup

use utility **pg_basebackup**

`pg_basebackup --username tomcat --pgdata /home/tomcat/backup-pgdata --label 'phisical_backup' --verbose`

outputs the next

```
pg_basebackup: initiating base backup, waiting for checkpoint to complete
pg_basebackup: checkpoint completed
pg_basebackup: write-ahead log start point: 0/4000028 on timeline 1
pg_basebackup: starting background WAL receiver
pg_basebackup: created temporary replication slot "pg_basebackup_1361"
pg_basebackup: write-ahead log end point: 0/4000100
pg_basebackup: waiting for background process to finish streaming ...
pg_basebackup: syncing data to disk ...
pg_basebackup: renaming backup_manifest.tmp to backup_manifest
pg_basebackup: base backup completed

```

and get the clone of pgdata directory

to verify the result of backup run **pg_verifybackup**

`localhost$ pg_verifybackup /home/tomcat/backup-pgdata`

```
backup successfully verified
```

