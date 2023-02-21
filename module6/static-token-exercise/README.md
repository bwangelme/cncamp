## 实验 - 为 api server 设置静态 token

- api-server 配置文件的位置

```shell
/etc/kubernetes/manifests/kube-apiserver.yaml
```

- 配置文件所做的修改

```yaml
command:
  - --token-auth-file=/etc/kubernetes/auth/static-token
volumeMounts:
  - mountPath: /etc/kubernetes/auth
    name: auth-files
    readOnly: true
volumes:
    - hostPath:
        path: /etc/kubernetes/auth
        type: DirectoryOrCreate
      name: auth-files
```

- 静态 token 文件 `static-token.csv` 的内容:

```csv
// token,user,uid,"groups"
secret-token,cncamp,1000,"group1,group2,group3"
```

- 使用 curl 访问 api-server

```shell
curl https://192.168.56.11:6443/api/v2/namespaces/default -H "Authorization: Bearer secret-token" -k
```

## 实验结果

可以看到 api-server 能够通过 token 正确地识别出用户 cncmap 来，说明认证已经通过了，报 403 的原因是鉴权没通过，cncmap 用户没有访问 namespaces 资源的权限。

```shell
ø> curl -k https://192.168.56.11:6443/api/v2/namespaces/default -H 'Authorization: Bearer secret-token'
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "namespaces \"default\" is forbidden: User \"cncamp\" cannot get resource \"namespaces\" in API group \"\" in the namespace \"default\"",
  "reason": "Forbidden",
  "details": {
    "name": "default",
    "kind": "namespaces"
  },
  "code": 403
}
```