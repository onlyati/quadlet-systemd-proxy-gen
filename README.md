# Systemd Quadlet proxy generator

This program generate the required socket and service files that allows casual
container activation, based on socket usage, with Podman Quadlet. For more
details check this
[post](https://thinkaboutit.tech/posts/2025-07-20-adhoc-containers-with-systemd-and-quadlet/).

## Example for usage

We have this Quadlet file, let it be called `nginx.container`:

```ini
[Unit]
Description=Nginx web server to server

[Container]
Image=docker.io/nginxinc/nginx-unprivileged
AutoUpdate=registry
User=%U
PublishPort=127.0.0.1:8080:8080

# Other
UserNS=keep-id:uid=101,gid=101

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5

```

After run the command, we got response:

```bash
$ quadlet-systemd-proxy-gen --quadlet nginx.container --ip 10.0.0.1
verify parameters:
- ip: 10.0.0.1
- port: 0
- container: nginx.container
- quadletIP: 127.0.0.1
- quadletPort: 0
creating socket and proxy files for ports: [8080]
generate file: /home/ati/.config/systemd/user/nginx-proxy-8080.socket
generate file: /home/ati/.config/systemd/user/nginx-proxy-8080.service

Post processing:
================
1. execute following commands to activate the generated data:
   systemctl --user daemon-reload
2. activate sockets
   be assume that [Unit] part contains the following in container files:
     nginx.container -> BindsTo=nginx-proxy-8080.service
     systemctl --user daemon-reload
   execute command
     systemctl --user enable --now nginx-proxy-8080.socket
```

Act as it is suggested, add the `BindsTo` line for Quadlet:

```ini
[Unit]
Description=Nginx web server to server
BindsTo=nginx-proxy-8080.service

[Container]
Image=docker.io/nginxinc/nginx-unprivileged
AutoUpdate=registry
User=%U
PublishPort=127.0.0.1:8080:8080

# Other
UserNS=keep-id:uid=101,gid=101

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5

```

Then execute commands accordingly:

```bash
$ systemctl --user daemon-reload
# Only systemd listening
$ sudo netstat -plnt | grep 8080
tcp        0      0 10.0.0.1:8080           0.0.0.0:*               LISTEN      1648/systemd
# Try to get request
$ curl 10.0.0.1:8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
# Container also listening on 8080 port
$ sudo netstat -plnt | grep 8080
tcp        0      0 10.0.0.1:8080           0.0.0.0:*               LISTEN      1648/systemd
tcp        0      0 127.0.0.1:8080          0.0.0.0:*               LISTEN      21071/pasta
# Container has been started
$ podman ps
CONTAINER ID  IMAGE                                         COMMAND               CREATED         STATUS         PORTS                     NAMES
108922dcf907  docker.io/nginxinc/nginx-unprivileged:latest  nginx -g daemon o...  15 seconds ago  Up 15 seconds  127.0.0.1:8080->8080/tcp  systemd-nginx
# 30 seconds after no connection, container stop
$ sleep 30 && podman ps
CONTAINER ID  IMAGE       COMMAND     CREATED     STATUS      PORTS       NAMES
```

## Download

You can go to GitHub release page to download the pre-built binary for your
computer or you can install via go:

```bash
go install github.com/onlyati/quadlet-systemd-proxy-gen@latest
```
