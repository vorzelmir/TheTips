### Interact Podman and minikube Kubernetes Cluster with GitLab container registry.

Before starting have to be done some actions on the gitlab.com

- Create a private project __"demo-server"__ for instance with deployment target "Registry (package or container)"

- Add new Project access token, the token name __"tom"__, Developer role and **api**, **read_registry**, **write_registry** permissions and copy the password


#### Push podman image to GitLab container registry.

Login podman with the token user to GitLab registry and save the authentication file explicitly

```
localhost$ podman login -u tom  --authfile ~/.config/containers/demo-server/auth.json registry.gitlab.com
Password: 
Login Succeeded!

```
Tag image "backendserver" [from](https://github.com/vorzelmir/TheTips/tree/main/nginx-reverse-proxy) on local machine  

`localhost$ podman tag localhost/backendserver registry.gitlab.com/<gitlab-user-name>/demo-server:v1`

Push taged image to the registry.gitlab.com

```
`localhost$ podman push registry.gitlab.com/<gitlab-user-name>/demo-server:v1`
.............
Copying config 8c2b9aaac3 done  
Writing manifest to image destination

```

#### Pull image from GitLab registry to minikube cluster.

Encode content of authentication file

`localhost$ cat ~/.confit/containers/demo-server/auth.json | base64`

Copy the result and pass it as an unbroken line to demo-server-secret.yaml file after .dockerconfigjson: 

Create file demo-server-secret.yaml

```
apiVersion: v1
kind: Secret
metadata:
  name: demo-server-gitlab-secret
data:
  .dockerconfigjson: ewoJImF....l9Cgl9Cn0=
type: kubernetes.io/dockerconfigjson
```
Create kubectl secret

`localhost$ kubectl apply -f demo-server-secret.yaml`

Check the result

`localhost$ kubectl get secret demo-server-gitlab-secret -o json | jq`

Or to compare with auth.json file

`localhost$ kubectl get secret demo-server-gitlab-secret  --output="jsonpath={.data.\.dockerconfigjson}" | base64 --decode`

Create a new demo-server-pod.yaml file

```
apiVersion: v1
kind: Pod
metadata:
  name: demo-server
spec:
  containers:
  - name: demo-server-container
    image: registry.gitlab.com/<gitlab-user-name>/demo-server:v1
  imagePullSecrets:
  - name: demo-server-gitlab-secret
```

Run new pod

`localhost$ kubectl apply -f demo-server-pod.yaml`

Check the result

```
localhost$ kubectl get pods
NAME          READY   STATUS    RESTARTS   AGE
demo-server   1/1     Running   0          5s

```
To confirm the pulling image

```
localhost$ kubectl describe pods demo-server | grep Pulled
  Normal  Pulled     102s  kubelet            Successfully pulled image "registry.gitlab.com/<gitlab-user-name>/demo-server:v1" in 1.759s (1.759s including waiting). Image size: 17421786 bytes.

```
