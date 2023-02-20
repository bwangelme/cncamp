## 启动错误

![](https://passage-1253400711.cos.ap-beijing.myqcloud.com/2023-02-19-211954.png)

日志

```shell
{"level":"warn","ts":"2023-02-20T01:12:54.569Z","caller":"etcdmain/etcd.go:146","msg":"failed to start etcd","error":"cannot access data directory: mkdir /bitnami/etcd/data: permission denied"}
{"level":"fatal","ts":"2023-02-20T01:12:54.569Z","caller":"etcdmain/etcd.go:204","msg":"discovery failed","error":"cannot access data directory: mkdir /bitnami/etcd/data: permission denied","stacktrace":"go.etcd.io/etcd/server/v3/etcdmain.startEtcdOrProxyV2\n\t/go/src/go.etcd.io/etcd/release/etcd/server/etcdmain/etcd.go:204\ngo.etcd.io/etcd/server/v3/etcdmain.Main\n\t/go/src/go.etcd.io/etcd/release/etcd/server/etcdmain/main.go:40\nmain.main\n\t/go/src/go.etcd.io/etcd/release/etcd/server/main.go:32\nruntime.main\n\t/go/gos/go1.16.15/src/runtime/proc.go:225"}
```

## 启动 etcd 集群的过程

### 在宿主机上创建目录

- 在每个节点上都要执行以下命令

```shell
sudo mkdir etcd_pv
# etcd server 默认是以 1001 这个用户启动的
sudo chown 1001:1001 etcd_pv/
```

### 创建 pv 和 storage-class

```shell
k apply -f pv.yaml
k apply -f storage-class.yaml
```

### 创建 etcd 集群

```shell
helm install my-etcd bitnami/etcd --set global.storageClass=manual --set persistence.size=1Gi --set replicaCount=3
```

### 测试

- 启动 etcd 集群并登陆进去

```shell
kubectl run my-etcd-client --restart='Never' --image docker.io/bitnami/etcd:3.5.7-debian-11-r10 --env ROOT_PASSWORD=$(kubectl get secret --namespace qae my-etcd -o jsonpath="{.data.etcd-root-password}" | base64 -d) --env ETCDCTL_ENDPOINTS="my-etcd.qae.svc.cluster.local:2379" --namespace qae --command -- sleep infinity
kubectl exec --namespace qae -it my-etcd-client -- bash
```

- 访问 etcd 集群

```shell
# 查看所有成员
I have no name!@my-etcd-client:/opt/bitnami/etcd$ etcdctl member list
5f9a9a7f98dfae79, started, my-etcd-0, http://my-etcd-0.my-etcd-headless.qae.svc.cluster.local:2380, http://my-etcd-0.my-etcd-headless.qae.svc.cluster.local:2379,http://my-etcd.qae.svc.cluster.local:2379, false
6a22aa1b49cbdad9, started, my-etcd-2, http://my-etcd-2.my-etcd-headless.qae.svc.cluster.local:2380, http://my-etcd-2.my-etcd-headless.qae.svc.cluster.local:2379,http://my-etcd.qae.svc.cluster.local:2379, false
90c3c58e20e4dc04, started, my-etcd-1, http://my-etcd-1.my-etcd-headless.qae.svc.cluster.local:2380, http://my-etcd-1.my-etcd-headless.qae.svc.cluster.local:2379,http://my-etcd.qae.svc.cluster.local:2379, false

# 读写 key
I have no name!@my-etcd-client:/opt/bitnami/etcd$ etcdctl --user root:$ROOT_PASSWORD put /message Hello
OK
I have no name!@my-etcd-client:/opt/bitnami/etcd$ etcdctl --user root:$ROOT_PASSWORD get /message
/message
Hello

# 查看集群的健康状态
I have no name!@my-etcd-client:/opt/bitnami/etcd$ etcdctl --user root:$ROOT_PASSWORD endpoint health
my-etcd.qae.svc.cluster.local:2379 is healthy: successfully committed proposal: took = 2.706814ms
```

## 参考链接

- [29 | PV、PVC体系是不是多此一举？从本地持久化卷谈起](https://time.geekbang.org/column/article/44245)
