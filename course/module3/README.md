## Docker

利用了 Linux 内核提供的技术，Namespace, CGroup, Union FS

## Namepsace

- 查看 namespace, lsns
- 查看所有的网络 namespace

```shell
root@lazyubuntu:~# lsns -t net
        NS TYPE NPROCS   PID USER         NETNSID NSFS                           COMMAND
4026531840 net     392     1 root      unassigned                                /sbin/init splash
4026532835 net       1   998 root      unassigned                                /usr/libexec/accounts-daemon
4026532902 net       1  1346 rtkit     unassigned                                /usr/libexec/rtkit-daemon
4026532970 net       2 12705 root               0 /run/docker/netns/eabedcedddd7 /bin/sh -c /httpserver
4026533532 net      14  4207 xuyundong unassigned                                /opt/google/chrome/chrome --type=zygote --crashpad-handler-pid=4198 --enable-crash-reporter=eab7d57c-c4e0-47fd-9657-ccb590fe4f80, --change-stack-guard-
4026533589 net       1  4208 xuyundong unassigned                                /opt/google/chrome/nacl_helper
```

- 查看一个进程的所有 namespace id

```shell
root@lazyubuntu:~# ls -l /proc/12705/ns/
总用量 0
lrwxrwxrwx 1 root root 0  1月  5 09:56 cgroup -> 'cgroup:[4026533036]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 ipc -> 'ipc:[4026532968]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 mnt -> 'mnt:[4026532966]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 net -> 'net:[4026532970]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 pid -> 'pid:[4026532969]'
lrwxrwxrwx 1 root root 0  1月  5 09:59 pid_for_children -> 'pid:[4026532969]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 time -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0  1月  5 09:59 time_for_children -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0  1月  5 09:56 uts -> 'uts:[4026532967]'
```

- 进入进程 12705 的 network namespace 执行 `ip a` 命令，查看所有的网卡信息

```shell
root@lazyubuntu:~# nsenter -t 12705 -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
29: eth0@if30: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```

- 启动的进程新建一个 Network Namespace, 并执行 `sleep 60` 命令

```shell
unshare -fn sleep 60

# 进入这个进程所在的 network namespace 执行 ip a 命令，可以看到只有一个 lo 网卡
root@lazyubuntu:~# nsenter -t 13380 -n ip a
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
```

## CGroup