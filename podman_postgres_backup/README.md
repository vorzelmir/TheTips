
Using podman version

`localhost$ podman version`


```
podman version
Client:       Podman Engine
Version:      4.6.2
API Version:  4.6.2
Go Version:   go1.19.12
Built:        Mon Aug 28 22:38:31 2023
OS/Arch:      linux/amd64

```

after Red Hat Registry service account 


`localhost$ podman login registry.redhat.io`


```
Username: {REGISTRY-SERVICE-ACCOUNT-USERNAME}
Password: {REGISTRY-SERVICE-ACCOUNT-PASSWORD}
Login Succeeded!

```

run


`localhost$ podman pull registry.redhat.io/rhel8/postgresql-16:1-16.1717586546`


check existing images


`localhost$ podman images | grep postgresql`


```
registry.redhat.io/rhel8/postgresql-16               latest            e54df09913ef  2 weeks ago   511 MB

```

create volume for data

`localhost$ podman volume create vol-postgres`

check volume exists

`localhost$ sudo ls $HOME/.local/share/containers/storage/volumes/`

```
...
drwx------.  3 user user 4096 Aug 28 12:21 vol-postgres
...
```

got the dump .sql file [from here](https://neon.tech/docs/import/import-sample-data#peridic)

copy .sql to the volume directory

`localhost$ cp peridic_table.sql $HOME/.local/share/containers/volumes/vol-postgres/_data`

run container 

`localhost$ podman run --name con-postgres -d -e POSTGRESQL_ADMIN_PASSWORD=secret -v vol-postgres:/var/lib/pgsql/data:Z --userns=keep-id:uid=26 rhel8/postgresql-16:latest `

where uid=26 is default image user's id

some logs 

`localhost$ podman logs con-postgres`

```
Starting server...
2024-06-25 08:01:27.805 UTC [1] LOG:  redirecting log output to logging collector process
2024-06-25 08:01:27.805 UTC [1] HINT:  Future log output will appear in directory "log".

```

persistent data in the cluster directory

`localhost$ podman exec -it con-postgres ls -l /var/lib/pgsql/data`

obtaing

```
-rw-rwxr--+  1 postgres     1000 17272 Jun 24 20:40 periodic_table.sql
drwx------. 20 postgres postgres  4096 Jun 25 08:01 userdata

```

connect to the container

`localhost$ podman exec -it con-postgres /bin/bash`

and go to the cluster directory /var/lib/pgsql/data

`bash-4.4$ cd  `

run .sql file 

`bash-4.4$ psql -U postgres -f periodic_table.sql`

check the result

`bash-4.4$ psql -U postgres -d periodic`
```
periodic=# \dt
             List of relations
 Schema |      Name      | Type  |  Owner   
--------+----------------+-------+----------
 public | periodic_table | table | postgres
(1 row)

```

describe table

`bash-4.4$ \dS+ periodic_tabel`

```
                                           Table "public.periodic_table"
      Column       |  Type   | Collation | Nullable | Default | Storage  | Compression | Stats target | Description 
-------------------+---------+-----------+----------+---------+----------+-------------+--------------+-------------
 AtomicNumber      | integer |           | not null |         | plain    |             |              | 
 Element           | text    |           |          |         | extended |             |              | 
 Symbol            | text    |           |          |         | extended |             |              | 
 AtomicMass        | numeric |           |          |         | main     |             |              | 
 NumberOfNeutrons  | integer |           |          |         | plain    |             |              | 
 NumberOfProtons   | integer |           |          |         | plain    |             |              | 
 NumberOfElectrons | integer |           |          |         | plain    |             |              | 
 Period            | integer |           |          |         | plain    |             |              | 

...........................................

 NumberOfIsotopes  | integer |           |          |         | plain    |             |              | 
 Discoverer        | text    |           |          |         | extended |             |              | 
 Year              | integer |           |          |         | plain    |             |              | 
 SpecificHeat      | numeric |           |          |         | main     |             |              | 
 NumberOfShells    | integer |           |          |         | plain    |             |              | 
 NumberOfValence   | integer |           |          |         | plain    |             |              | 
Indexes:
    "periodic_table_pkey" PRIMARY KEY, btree ("AtomicNumber")
Access method: heap

```

exit from container and run *pg_dump* to dump the periodic database

`bash4-4$ exit`

`localhost$ podman exec -it con-postgres pg_dump -U postgres -d periodic -f periodic.dump.sql`

check the result inside volume vol-postgres

`localhost$ sudo ls -la $USER/.local/share/containers/volumes/vol-postgres/_data`

```
drwxrwx---.  3 user 100000  4096 Aug 28 12:27 .
drwx------.  3 user user     4096 Aug 28 12:21 ..
-rw-r--r--.  1 user 100026 17272 Aug 28 12:27 periodic.dump.sql
-rw-r-xr--.  1 user user    17272 Aug 28 12:20 periodic_table.sql
...
```
