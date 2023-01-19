Kubernetes 架构原则和对象设计
===

## 操作 etcd

```shell
export ETCDCTL_API=3
# 查看 / 下的所有key
k -n kube-system exec -it etcd-k8s-node1 -- etcdctl --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --cacert /etc/kubernetes/pki/etcd/ca.crt get --keys-only --prefix /
# 查看 key /registry/services/specs/default/kubernetes
k -n kube-system exec -it etcd-k8s-node1 -- etcdctl --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --cacert /etc/kubernetes/pki/etcd/ca.crt get --prefix /registry/services/specs/default/kubernetes
```

## 控制器协同工作原理

控制器操作 Deployment 的流程图

![](https://passage-1253400711.cos.ap-beijing.myqcloud.com/2023-01-19-120404.png)

```shell
# 查看 pod 的事件，可以看到第一件是将 pod assigned 到 k8s-node3 这件事是 default-scheduler 做的 
ø> k describe pod nginx-deployment-8f458dc5b-xlccw | tail
Events:
  Type    Reason     Age    From               Message
  ----    ------     ----   ----               -------
  Normal  Scheduled  9m28s  default-scheduler  Successfully assigned default/nginx-deployment-8f458dc5b-xlccw to k8s-node3
  Normal  Pulling    9m27s  kubelet            Pulling image "nginx"
  Normal  Pulled     9m25s  kubelet            Successfully pulled image "nginx" in 2.549180878s
  Normal  Created    9m25s  kubelet            Created container nginx
  Normal  Started    9m25s  kubelet            Started container nginx

```