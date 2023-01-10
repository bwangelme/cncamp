OOM 的练习
===

这是一个测试 cgroup memory 子系统的 demo, 程序的主要功能是每分钟分配 10Mb 的内存，而且这部分内存不会被 golang GC 回收。

执行 make build, 会生成二进制文件 bin/malloc


## 操作步骤

- 创建 memorydemo 控制器

```shell
   35  mkdir memorydemo
   36  cd memorydemo/
```

- 设置控制器的内存限制

```shell
# 这意味着此控制器设置的最大内存是 10Mb
   38  echo 104960000 > memory.limit_in_bytes
```

- 启动 malloc 进程, 并查看进程号

```shell
./malloc
```

```shell
# 此命令能够查看进程使用的内存以及进程号, RSS 查看实际使用的内存
watch 'ps -aux | grep malloc | grep -v grep '
```

- 将 malloc 进程加入到 memorydemo 控制器中

__注意:__ 一定要先设置控制器内存限制大小(`memory.limit_in_bytes`), 再向控制器中添加进程

如果控制器已经有进程了，再修改内存限制，会出错:

```shell
root@dockervbox:/sys/fs/cgroup/memory/memorydemo# echo 104960000 > memory.limit_in_bytes
-bash: echo: write error: Device or resource busy
```

```shell
echo 3622 > cgroup.procs
```

- 等一分钟后，可以看到进程被杀掉了

```shell
vagrant@dockervbox:/vagrant/bin$ ./malloc
Allocating 100Mb memory, raw memory is 104960000
Allocating 200Mb memory, raw memory is 209920000
Killed
```

- 在 kernel 日志中可以看到 oom 的记录

```shell
root@dockervbox:~# egrep -i -r 'Out of memory' /var/log/
/var/log/kern.log:Jan 10 01:27:05 dockervbox kernel: [ 1683.797227] Memory cgroup out of memory: Killed process 3985 (malloc) total-vm:1211812kB, anon-rss:205140kB, file-rss:872kB, shmem-rss:0kB, UID:1000 pgtables:524kB oom_score_adj:0
/var/log/syslog:Jan 10 01:27:05 dockervbox kernel: [ 1683.797227] Memory cgroup out of memory: Killed process 3985 (malloc) total-vm:1211812kB, anon-rss:205140kB, file-rss:872kB, shmem-rss:0kB, UID:1000 pgtables:524kB oom_score_adj:0
Binary file /var/log/journal/e0e8acfb951d475d9f1b4ed29f1fa00c/system.journal matches

root@dockervbox:~# dmesg -k | grep 'out of memory'
[ 1683.797227] Memory cgroup out of memory: Killed process 3985 (malloc) total-vm:1211812kB, anon-rss:205140kB, file-rss:872kB, shmem-rss:0kB, UID:1000 pgtables:524kB oom_score_adj:0
```