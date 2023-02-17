## 生成证书

- clone etcd 的代码，并进入 `etcd/hack/tls-setup` 目录
- 编辑 csr 文件 `v config/req-csr.json`

```
{
  "CN": "etcd",
  "hosts": [
    "localhost",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "O": "GeekBang",
      "OU": "CNCamp",
      "L": "BeiJing"
    }
  ]
}
```

- 生成证书 

```shell
infra0=127.0.0.1 infra1=127.0.0.1 infra2=127.0.0.1 make
```

- 将证书放到当前目录的 certs 目录中

## 启动集群

```shell
mkdir data log
```

`./start-all.sh` 脚步可以启动集群，我将数据文件放到了当前目录的 `data` 中，日志放到了 `log` 中

## 查看集群

```shell
etcdctl --endpoints https://127.0.0.1:3379 --cert certs/127.0.0.1.pem --key certs/127.0.0.1-key.pem --cacert certs/ca.pem member list

ø> etcdctl --endpoints https://127.0.0.1:3379 --cert certs/127.0.0.1.pem --key certs/127.0.0.1-key.pem --cacert certs/ca.pem member list
1701f7e3861531d4, started, infra0, https://127.0.0.1:3380, https://127.0.0.1:3379, false
6a58b5afdcebd95d, started, infra1, https://127.0.0.1:4380, https://127.0.0.1:4379, false
84a1a2f39cda4029, started, infra2, https://127.0.0.1:5380, https://127.0.0.1:5379, false
```

## 写入数据

```shell
etcdctl --endpoints https://127.0.0.1:3379 --cert certs/127.0.0.1.pem --key certs/127.0.0.1-key.pem --cacert certs/ca.pem put a b
etcdctl --endpoints https://127.0.0.1:3379 --cert certs/127.0.0.1.pem --key certs/127.0.0.1-key.pem --cacert certs/ca.pem put c d
```

## 备份数据

```shell
etcdctl --endpoints https://127.0.0.1:3379 --cert certs/127.0.0.1.pem --key certs/127.0.0.1-key.pem --cacert certs/ca.pem snapshot save snapshot.db
```

## 删除数据，停止集群

```shell
./stop-all.sh
rm -r data
```

## 恢复数据

```shell
./restore.sh
```

## 重启集群

- __NOTE__: 因为 restore.sh 中在恢复数据的时候，已经指定了 `--initial-cluster`, `--initial-cluster-token`, `--initial-advertise-peer-urls` 参数了，所以再次启动集群的时候，就不需要再指定了，这些数据已经在 etcd 中存储了。

```shell
./restart-all.sh
```

