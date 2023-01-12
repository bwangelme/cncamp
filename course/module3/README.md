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

### cgroup v1 提供的 cpu 接口

- `cpu.shares` 进程占用 cpu 时间的相对比例，如果只有一个进程，那么会全部占用，如果有多个进程，根据设置的值的大小来分配 CPU 时间。和 k8s.resources.cpu 是相同的概念
- `cpu.cfs_period_us` 设置时间周期的长度，单位是 us (微秒)
- `cpu.cfs_quota_us` 设置当前 cgroup 在 `cfs_period_us` 设置的周期内，最多能运行多长时间，单位是 us.
- `cpu.stat` cgroup 内进程使用的 CPU 时间统计
- `nr_periods` 经过 `cfs_period_us` 设置的周期数量
- `nr_throttled` 在经过的周期内，有多少次是因为进程在指定的时间周期内用光了配额时间而收到限制
- `throttled_time` cgroup 中的进程被限制使用的 CPU 的总用时，单位是 ns(纳秒)

### cgroup v2 提供的 CPU 接口

- `cgroup.procs` 设置当前 cgroup 管理的进程 ID
- `cpu.max` 格式是 `$MAX $PERIOD`， 类似于 `cfs_period_us cfs_quota_us`, 设置每个周期的时间长度，以及周期内可以使用的时间，单位是 us (microseconds)

### CFS 完全公平调度

- CFS (Completely Fail Scheduler), 完全公平调度器。
- CFS 的思想是维护为任务提供处理器时间方面的平衡，这意味着给进程分配相当数量的处理器。
- CFS 中的 vruntime 概念:
  - vruntime = 实际运行时间 * 1024 / 进程权重
  - CFS 记录每个进程的运行时间，根据不同的权重，记录的运行时间也会不同，权重越大，vruntime 越小
- CFS 维护一颗红黑树，将进程按照 vruntime 进行排列，vruntime 最小的在最左边，每次调度时，选取最左边的进程执行
- 红黑树的特点:
  - 自平衡，树上没有一条路径会比其他的路径长两倍
  - O(log n) 的时间复杂度，能够在树上快速高效地插入或删除进程

## Overlayfs 实验

```shell
root@dockervbox:~# mkdir overlayfs
root@dockervbox:~# cd overlayfs/
# 创建四个目录，lower,upper 的文件将要被 merge 到 merged 目录中，work 是 overlayfs 的工作目录
root@dockervbox:~/overlayfs# mkdir upper lower merged work
root@dockervbox:~/overlayfs# echo "from lower" > lower/in_lower.txt
root@dockervbox:~/overlayfs# echo "from upper" > upper/in_upper.txt
root@dockervbox:~/overlayfs# echo "from lower" > lower/in_both.txt
root@dockervbox:~/overlayfs# echo "from upper" > upper/in_both.txt
root@dockervbox:~/overlayfs# tree .
.
├── lower
│   ├── in_both.txt
│   └── in_lower.txt
├── merged
├── upper
│   ├── in_both.txt
│   └── in_upper.txt
└── work

4 directories, 4 files
# 这条命令将一个 overlayfs 挂载到 merged 目录中
root@dockervbox:~/overlayfs# mount -t overlay overlay -o lowerdir=`pwd`/lower,upperdir=`pwd`/upper,workdir=`pwd`/work `pwd`/merged
root@dockervbox:~/overlayfs# mount | grep -i overlay
overlay on /root/overlayfs/merged type overlay (rw,relatime,lowerdir=/root/overlayfs/lower,upperdir=/root/overlayfs/upper,workdir=/root/overlayfs/work,xino=off)
root@dockervbox:~/overlayfs# tree merged/
merged/
├── in_both.txt
├── in_lower.txt
└── in_upper.txt

0 directories, 3 files
# 可以看到，upper 和 lower 目录的文件被合并到了 merged 目录中，文件名相同的文件，使用的是 upper 目录的文件
root@dockervbox:~/overlayfs# cat merged/in_both.txt
from upper
root@dockervbox:~/overlayfs# cat merged/in_lower.txt
from lower
root@dockervbox:~/overlayfs# cat merged/in_upper.txt
from upper
root@dockervbox:~/overlayfs#
```

## Docker 架构

![](https://passage-1253400711.cos.ap-beijing.myqcloud.com/2023-01-12-082246.png)

docker 的架构如上图所示，containerd 通过 shim 启动了容器子进程之后，它不直接管理，shim 进程是 systemd-init 的子进程，所以 containerd 自身更加轻量，没有子进程也易于重启。

![](https://passage-1253400711.cos.ap-beijing.myqcloud.com/2023-01-12-082557.png)

## docker 使用 none 网络模式从 0 构建网络的实验

```shell
# 查看当前主机上的网桥设备，默认有一个 docker0
root@dockervbox:~# brctl show
bridge name     bridge id               STP enabled     interfaces
docker0         8000.02421f335f1f       no

# 创建 /var/run/netns 目录，ip netns list 查看的就是此目录中的网络 namespace
root@dockervbox:~# mkdir -p /var/run/netns

# 删除所有的网络 ns
root@dockervbox:~# find -L /var/run/netns -type l -delete

# 使用 none 网络模式启动一个 nginx 容器
root@dockervbox:~# docker run --network=none -d nginx
Unable to find image 'nginx:latest' locally
latest: Pulling from library/nginx
8740c948ffd4: Pull complete
d2c0556a17c5: Pull complete
c8b9881f2c6a: Pull complete
693c3ffa8f43: Pull complete
8316c5e80e6d: Pull complete
b2fe3577faa4: Pull complete
Digest: sha256:b8f2383a95879e1ae064940d9a200f67a6c79e710ed82ac42263397367e7cc4e
Status: Downloaded newer image for nginx:latest
305efb91b1ba38d2ebb9b64c499fb84bcde9ff6b1f26eadebd08fbe6fc8898b6

root@dockervbox:~# docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS          PORTS     NAMES
305efb91b1ba   nginx     "/docker-entrypoint.…"   13 seconds ago   Up 12 seconds             nervous_herschel

# 启动成功后，记录此容器的 pid 到 pid环境变量
root@dockervbox:~# docker inspect 305efb91b1ba -i pid
unknown shorthand flag: 'i' in -i
See 'docker inspect --help'.
root@dockervbox:~# docker inspect 305efb91b1ba | grep -i pid
            "Pid": 4955,
            "PidMode": "",
            "PidsLimit": null,
root@dockervbox:~# export pid=4955
root@dockervbox:~# echo $pid
4955

# 进入 $pid 进程所在的网络 ns, 查看所有的网络设备，可以看到只有一个 lo 设备，没有其他的网络设备
root@dockervbox:~# nsenter -t $pid -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
       
# 软链接 $pid 的网络 ns 到 /var/run/netns/ 中, 创建 ip netns 可以操作的网络 ns
root@dockervbox:~# ln -s "/proc/$pid/ns/net" /var/run/netns/$pid
root@dockervbox:~# ls /var/run/netns/4955
/var/run/netns/4955
root@dockervbox:~# ip netns list
4955

# 创建一对 veth 设备，两端设备分别叫做 A 和 B
root@dockervbox:~# ip link add A type veth peer name B
# 将 A 连接到 docker0 网桥设备上
root@dockervbox:~# brctl addif docker0 A
# 将 A 启动
root@dockervbox:~# ip link set A up
# 在主机网络 ns 中查看所有的设备，可以看到 A@B 已经启动了
root@dockervbox:~# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: enp0s3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether 02:b7:1d:9c:e0:75 brd ff:ff:ff:ff:ff:ff
    inet 10.0.2.15/24 brd 10.0.2.255 scope global dynamic enp0s3
       valid_lft 82712sec preferred_lft 82712sec
    inet6 fe80::b7:1dff:fe9c:e075/64 scope link
       valid_lft forever preferred_lft forever
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default
    link/ether 02:42:1f:33:5f:1f brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever
4: B@A: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN group default qlen 1000
    link/ether a6:a1:28:a0:94:bf brd ff:ff:ff:ff:ff:ff
    
# 看起来没有完全启动成功，还有一个 M-DOWN 的状态
5: A@B: <NO-CARRIER,BROADCAST,MULTICAST,UP,M-DOWN> mtu 1500 qdisc noqueue master docker0 state LOWERLAYERDOWN group default qlen 1000
    link/ether 16:05:7b:c1:c3:86 brd ff:ff:ff:ff:ff:ff
# 将容器的目标 ip, 子网掩码，网关保存到环境变量中
root@dockervbox:~# SETIP=172.17.0.10
root@dockervbox:~# SETMASK=16
root@dockervbox:~# GATEWAY=172.17.0.1

# 将 veth B 放到 $pid 网络 ns 中
root@dockervbox:~# ip link set B netns $pid

# 将 veth B 在 $pid ns 中的名字改成 eth0
root@dockervbox:~# ip netns exec $pid ip link set dev B name eth0

# 启动 $pid ns 中的 eth0
root@dockervbox:~# ip netns exec $pid ip link set eth0 up

# 查看 $pid ns 中的所有进程，此时可以看到 eth0 已经启动了
# 因为 A 已经 连接到了网桥上，所以直接启动成功，出现了 LOWER_UP 的状态
root@dockervbox:~# ip netns exec $pid ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether a6:a1:28:a0:94:bf brd ff:ff:ff:ff:ff:ff link-netnsid 0
    
# 设置 $pid ns 中的 eth0 设备的 ip 和 掩码
root@dockervbox:~# ip netns exec $pid ip addr add $SETIP/$SETMASK dev eth0
# 设置 $pid ns 中的 eth0 设备的默认路由网关
root@dockervbox:~# ip netns exec $pid ip route add default via $GATEWAY
# 此时访问 nginx 容器的 ip, 发现已经能够访问通了
root@dockervbox:~# curl $SETIP
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

# 查看 $pid ns 中的所有设备，可以看到 eth0 已经启动成功，并且有了 ip 和子网掩码
root@dockervbox:~# nsenter -t $pid -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether a6:a1:28:a0:94:bf brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.10/16 scope global eth0
       valid_lft forever preferred_lft forever
       
# 查看 $pid ns 中的路由表, 可以看到 eth0 设备的默认路由地址是 
root@dockervbox:~# nsenter -t $pid -n ip r 172.17.0.1
default via 172.17.0.1 dev eth0
172.17.0.0/16 dev eth0 proto kernel scope link src 172.17.0.10

# 查看宿主机上的网络设备，可以看到 veth A 的名字变成了 A@if4, 并且状态已经由 M-DOWN 变成了 LOWER_UP
root@dockervbox:~# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: enp0s3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    link/ether 02:b7:1d:9c:e0:75 brd ff:ff:ff:ff:ff:ff
    inet 10.0.2.15/24 brd 10.0.2.255 scope global dynamic enp0s3
       valid_lft 81989sec preferred_lft 81989sec
    inet6 fe80::b7:1dff:fe9c:e075/64 scope link
       valid_lft forever preferred_lft forever
3: docker0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether 02:42:1f:33:5f:1f brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:1fff:fe33:5f1f/64 scope link
       valid_lft forever preferred_lft forever
5: A@if4: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue master docker0 state UP group default qlen 1000
    link/ether 16:05:7b:c1:c3:86 brd ff:ff:ff:ff:ff:ff link-netns 4955
    inet6 fe80::1405:7bff:fec1:c386/64 scope link
       valid_lft forever preferred_lft forever
```