## 启动容器

```shell
vagrant@dockervbox:~$ docker run --rm -it -p 8080:8080 bwangel/cncamp_http_server:v1.0
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /healthz                  --> main.Healthz (4 handlers)
[GIN-debug] GET    /                         --> main.Home (4 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8080
```

## 查看本镜像启动后，容器的 ip 配置

```shell
vagrant@dockervbox:~$ docker ps
CONTAINER ID   IMAGE                             COMMAND                  CREATED         STATUS         PORTS                                       NAMES
49cde08d1d15   bwangel/cncamp_http_server:v1.0   "/bin/sh -c /httpser…"   5 minutes ago   Up 5 minutes   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   reverent_chebyshev
vagrant@dockervbox:~$ docker inspect 49cde08d1d15 | grep -i pid
            "Pid": 11854,
            "PidMode": "",
            "PidsLimit": null,
vagrant@dockervbox:~$ ps aufx | grep 11854
vagrant    12645  0.0  0.0   8160   720 pts/1    S+   01:28   0:00              \_ grep --color=auto 11854
root       11854  0.0  0.0   2880   940 pts/0    Ss+  01:23   0:00  \_ /bin/sh -c /httpserver
vagrant@dockervbox:~$ sudo nsenter -t 11854 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
6: eth0@if7: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```
