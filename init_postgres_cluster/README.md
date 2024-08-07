### run postgresql server

operation system is

`server$ hostnamectl | grep Operating`

```
Operating System: Fedora Linux 39 (Server Edition)

```

create user tomcat and add his to postgres group

`server# useradd --groups postgres --create-home tomcat`

add to users from postgres group permission to create .lock files

`server# chmod g+w /var/run/postgresql`


create postgresql cluster

`server$ pg_ctl init --pgdata=/home/tomcat/pgdata `

```
...
initdb: hint: You can change this by editing pg_hba.conf or using the option -A, or --auth-local and --auth-host, the next time you run initdb.

Success. You can now start the database server using:

    /usr/bin/pg_ctl -D /home/tomcat/pgdata -l logfile start
```

start postgresql server

`server$ pg_ctl --pgdata=/home/tomcat/pgdata --log=log/pgctl.log start`

```
waiting for server to start.... done
server started
```

assign PGDATA explicitly

`server$ export PGDATA=/home/tomcat/pgdata`

and check the result

`server$ pg_ctl status`

```
pg_ctl: server is running (PID: 1111)
/usr/bin/postgres "-D" "pgdata"
```

connect to cluster and get the list of default databases for examples

`server$ psql -l`

```
                                              List of databases
   Name    | Owner  | Encoding |   Collate   |    Ctype    | ICU Locale | Locale Provider | Access privileges 
-----------+--------+----------+-------------+-------------+------------+-----------------+-------------------
 postgres  | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | 
 template0 | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | =c/tomcat        +
           |        |          |             |             |            |                 | tomcat=CTc/tomcat
 template1 | tomcat | UTF8     | en_US.UTF-8 | en_US.UTF-8 |            | libc            | =c/tomcat        +
           |        |          |             |             |            |                 | tomcat=CTc/tomcat
(3 rows)
```


