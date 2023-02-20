## 启动错误

![](https://passage-1253400711.cos.ap-beijing.myqcloud.com/2023-02-19-211954.png)

日志

```shell
{"level":"warn","ts":"2023-02-20T01:12:54.569Z","caller":"etcdmain/etcd.go:146","msg":"failed to start etcd","error":"cannot access data directory: mkdir /bitnami/etcd/data: permission denied"}
{"level":"fatal","ts":"2023-02-20T01:12:54.569Z","caller":"etcdmain/etcd.go:204","msg":"discovery failed","error":"cannot access data directory: mkdir /bitnami/etcd/data: permission denied","stacktrace":"go.etcd.io/etcd/server/v3/etcdmain.startEtcdOrProxyV2\n\t/go/src/go.etcd.io/etcd/release/etcd/server/etcdmain/etcd.go:204\ngo.etcd.io/etcd/server/v3/etcdmain.Main\n\t/go/src/go.etcd.io/etcd/release/etcd/server/etcdmain/main.go:40\nmain.main\n\t/go/src/go.etcd.io/etcd/release/etcd/server/main.go:32\nruntime.main\n\t/go/gos/go1.16.15/src/runtime/proc.go:225"}
```