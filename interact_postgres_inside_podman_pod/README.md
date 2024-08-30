#### create podman pod

podman version

`localhost$ podman version`

```
Client:       Podman Engine
Version:      4.6.2
API Version:  4.6.2
Go Version:   go1.19.12
Built:        Mon Aug 28 22:38:31 2023
OS/Arch:      linux/amd64
```

create pod

`localhost$ podman pod create --publish 5432:5432 --name pod-post-go`

create postgres container inside of the pod from official image

`localhost$ podman create --pod pod-post-go -e POSTGRES_PASSWORD=secret --name con-postgres postgres`

make image from Golang main.go file to get the version of Postgres container

`localhost$ buildah build -f Containerfile --tag img-golang`


```
...
Successfully tagged localhost/img-golang:latest
...
```
create go container from this image

`localhost$ podman create --pod pod-post-go --name con-golang localhost/img-golang`

watch the result

`localhost$ podman ps -a`

```
CONTAINER ID  IMAGE                                    COMMAND     CREATED         STATUS      PORTS                   NAMES
a0fa12581ba8  localhost/podman-pause:4.6.2-1693251511              56 seconds ago  Created     0.0.0.0:5432->5432/tcp  30a7554a0276-infra
60e0a201681a  docker.io/library/postgres:latest        postgres    35 seconds ago  Created     0.0.0.0:5432->5432/tcp  con-postgres
6536ba42c7d7  localhost/img-golang:latest                          14 seconds ago  Created     0.0.0.0:5432->5432/tcp  con-golang

```

containers created but not running

run the pod 

`localhost$ podman pod start pod-post-go`

check the result of running pod

`localhost$ podman ps -a`

```
CONTAINER ID  IMAGE                                    COMMAND     CREATED        STATUS                    PORTS                   NAMES
a0fa12581ba8  localhost/podman-pause:4.6.2-1693251511              3 minutes ago  Up 11 seconds             0.0.0.0:5432->5432/tcp  30a7554a0276-infra
60e0a201681a  docker.io/library/postgres:latest        postgres    3 minutes ago  Up 11 seconds             0.0.0.0:5432->5432/tcp  con-postgres
6536ba42c7d7  localhost/img-golang:latest                          2 minutes ago  Exited (0) 8 seconds ago  0.0.0.0:5432->5432/tcp  con-golang

```

get the  postgres version

`localhost$ podman logs con-golang`

```
PostgreSQL 16.3 (Debian 16.3-1.pgdg120+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bi
```
