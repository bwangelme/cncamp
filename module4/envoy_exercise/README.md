## 说明

文件名|作用
---|---
envoy-configmap.yaml|envoy 配置文件的 configmap, 最终会被挂载进 pod 中
envoy_deploy.yaml| envoy deployment 的配置文件
envoy_svc.yaml | envoy service 的配置

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