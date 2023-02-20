etcd 命令
===

```shell
alias etcl='etcdctl --endpoints https://127.0.0.1:2379 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key'

etcl get --prefix --keys-only /

# namespace 的 http endpoint 是 /api/v1/namespaces/qae
# namespace 在 etcd 中的存储路径是
etcl get /registry/namespaces/qae
```

