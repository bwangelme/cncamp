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