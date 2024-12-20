## basic nginx configuration in RedHat-9.4

check available nginx version in repository

`server# yum module list nginx`

```
Red Hat Enterprise Linux 9 for x86_64 - AppStream (RPMs)
Name                          Stream                           Profiles                          Summary                               
nginx                         1.22                             common [d]                        nginx webserver                       
nginx                         1.24                             common                            nginx webserver  

```
enable last version and install it

`server# yum module enable nginx:1.24`

`server# yum install nginx`

start nginx server

`server# systemctl enable nginx --now`

configure the firewall: make zone "home" active and enable http, https services inside it

`server# firewall-cmd --set-default-zone=home`

`server# firewall-cmd --zone=home --add-service={http,https}`

`server# firewall-cmd --zone=home --add-service={http,https} --permanent`

### run two virtual domains first.com and second.net

create two directories in /var/www

`server# mkdir -p /var/www/{first.com,second.net}`

add content inside

`server# echo "hello from first.com" > /var/www/first.com/index.html`

`server# echo "hello from second.net" > /var/www/second.net/index.html`

check the selinux context type, it must be httpd_sys_content_t. if does not change it

`server# semanage fcontext -a -t httpd_sys_content_t "/var/www/{first.com,second.net}(/.*)?"`

`server# restorecon -R "/var/www/{first.com,second.net}"`

append two new server{} contexts to the /etc/nginx/nginx.conf inside http{} context

```
#new context to add a first.com server

server {
    domain_name first.com;
    root  /var/www/first.com;
}

#new context to add second.net server

server {
    domain_name second.net;
    root /var/www/second.net;
}
```

in the client machine add first.com and second.net with ip addresses of the server machine to the /etc/hosts

check the result

`client$ curl first.com`

or

`client$ curl second.net`

### add context location

create directorie and hello.txt file inside

`server# mkdir /var/www/first.com/docs`

`server# echo "hello from Tom's docs" > /var/www/first.com/greet.txt`

edit /etc/nginx/nginx.conf

adding location directives in the server context

```
server {
         server_name first.com;
         root /var/www/first.com;
         location  = /docs {
             root /docs;
             try_files $uri = 400;
         }   
     }
```

send request from host machine

`client$ curl first.com/docs/greet.txt`

get Tom greeting from /var/www/first.com/docs/greet.txt


### enable TLS

create dir in home directories for keys, requests and certificates

`server$ mkdir keys`

create private key Certificate Authorities for domain second.net

`server# openssl genpkey -algorithm rsa -pkeyopt rsa_keygen_bits:2048 -out keys/second.net.key`

create new request signed with private key

`server# openssl req -new -key keys/second.net.key -out keys/second.net.csr`

get certificate based on request

`server# openssl x509 -req -days 365 -in keys/second.net.csr -signkey keys/second.net.key -out keys/second.net.crt`

copy key and certificate to the /etc/pki/tls

`server# cp keys/second.net.key /etc/pki/tls/private/`

`server# cp keys/second.net.crt /etc/pki/tls/certs/`

and check the result

```
client$ curl https://second.net
curl: (60) SSL certificate problem: self-signed certificate
More details here: https://curl.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the web page mentioned above.

```

to ignore self-signed error and get more info

`client$ curl --ignore -vv https://second.net`

