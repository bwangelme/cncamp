## 测试

- 在 k8s 节点上，通过 svc ip 可以访问 envoy 容器

```shell
ø> k get svc                                                                                                                                                                                                           09:43:08 (02-01)
NAME        TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
envoy-svc   ClusterIP   10.111.186.177   <none>        8080/TCP   4s

vagrant@k8s-node2:~$ curl http://10.111.186.177:8080/get
{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "www.httpbin.org",
    "User-Agent": "curl/7.68.0",
    "X-Amzn-Trace-Id": "Root=1-63d9c3dd-3173f87522d6b1a24e472294",
    "X-Envoy-Expected-Rq-Timeout-Ms": "15000"
  },
  "origin": "123.117.183.217",
  "url": "https://www.httpbin.org/get"
}
```