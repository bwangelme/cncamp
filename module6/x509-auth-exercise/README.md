## 创建 key 和 csr 文件

```shell
openssl genrsa --out myuser.key 2048
openssl req -new -key myuser.key -out myuser.csr
```

```shell
ø> openssl req -new --key myuser.key --out myuser.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:BeiJing
Locality Name (eg, city) []:BeiJing
Organization Name (eg, company) [Internet Widgits Pty Ltd]:cncmap
Organizational Unit Name (eg, section) []:cncmap
# 注意: 这里填入的 Common Name, 就是在 k8s 中定义的用户名
Common Name (e.g. server FQDN or YOUR name) []:cncamp
Email Address []:cncmap@gmail.com

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```

## 给 k8s 发送签发证书的请求

```shell
k applf -f csr.yaml
```

## k8s 签发证书

```shell
k certificate approve myuser
```

## 通过证书访问 k8s

```shell
ø> k get csr myuser -o jsonpath='{.status.certificate}'| base64 -d > myuser.crt
ø> k config set-credentials myuser --client-key=myuser.key --client-certificate=myuser.crt --embed-certs=true
User "myuser" set.
```

```shell
# myuser 可以通过认证，但是没有访问权限
ø> k get pod --user myuser
Error from server (Forbidden): pods is forbidden: User "cncamp" cannot list resource "pods" in API group "" in the namespace "qae"
```

## 创建 role 和 role binding, 让用户可以 list pod

```shell
ø> k create role developer --verb=get --verb=list --verb=update --verb=delete --resource=pods
role.rbac.authorization.k8s.io/developer created
# cncamp 用户名是在生成 csr 文件时，填入的 CN
ø> k create rolebinding developer-binding-myuser --role=developer --user=cncamp
rolebinding.rbac.authorization.k8s.io/developer-binding-myuser created
```