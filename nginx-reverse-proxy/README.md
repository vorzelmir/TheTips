### Configure nginx as a reverse proxy on RHEL-9.4

#### some steps in client machine

ping proxy server

```
ping -c 1 proxy-server.com
PING proxy-server.com (ip-address proxy-server) 56(84) bytes of data.
64 bytes from ip-address proxy-server: icmp_seq=1 ttl=64 time=0.596 ms

```
ping to proxy-server is going

ping backend server

```
ping -c 1 backend-server.com 
PING backend-server.com (ip-address backend-server.com) 56(84) bytes of data.
From ip-address client machine icmp_seq=1 Destination Host Unreachable
```
ping does not go through because backend-server and client are in the different subnets


#### configure proxy server

add next content to the /etc/nginx/nginx.conf

```

http {

..................
    server {
        server_name backend-server.com;
        location / {
            proxy_pass http://backend-server.com:8888;
        }
    }
}
```


allow SELinux to pass traffic through proxy server 

`proxy# setsebool -PV httpd_can_connect_network=on`

#### backend server

backend-server is a simple server written on Go with only two options: return current date/time and the name of Operation System

backend-server is disabled by default and to make it running when proxy-server sends a request it needs to make sysetmd unit.service

some steps on this task are needed:

  1. Create /etc/systemd/system/backend-server.service file 
  ```
   [Unit]
   Description=backend-server
   After=network.target
   
   [Service]
   User=root
   Group=root
   Type=simple
   Restart=on-failure
   RestartSec=5s
   ExecStart=/home/User/bin/backend-server
   
   [Install]
   WantedBy=multi-user.target

  ```

  2. Run the unit file now and enable the next starting

  `backend-server# systemctl enable --now backend-server.service`
  
  3. Check the status of backend-server.service

   `backend-server$ systemctl status backend-server`

   if not active check SELinux warnings

   `backend-server# ausearch -m avc`

   SELinux may not get permission systemd to run the unit file if it is not in the default executable directories,
   so check server's $PATH environment and add a binary backend-server file to the /home/User/bin directory for instance

   repeat step 2

  4. Configure the firewall to allow tcp requests on 8888 port

  `backend-server# firewall-cmd --add-port=8888/tcp`
  
  `backend-server# firewall-cmd --add-port=8888/tcp --permanent`


#### client machine

send http requests to backend-server and get results

```
curl backend-server.com/time
Sunday, 19-Jan-25 09:09:43 EST

```

```
curl backend-server.com/os
Operational System ="Fedora Linux 40 (Server Edition)"

```
