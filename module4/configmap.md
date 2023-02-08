```shell
# 将整个文件的内容创建成一个 key
$ kubectl create configmap game-config --from-file=game.properties
ø> k get configmaps game-config -o yaml
apiVersion: v1
data:
  game.properties: |-
    enemies=aliens
    bbb
    #aaa
    lives=3
    enemies.cheat=true
    enemies.cheat.level=noGoodRotten
    secret.code.passphrase=UUDDLRLRBABAS
    secret.code.allowed=true
kind: ConfigMap
metadata:
  name: game-config
  namespace: qae

# 解析 game.properties 文件，将里面的键值对解析出来
# 忽略了加井号的注释行
$ kubectl create configmap game-env-config --from-env-file=game.properties

ø> k get configmaps game-env-config -o yaml
apiVersion: v1
data:
  bbb: ""
  enemies: aliens
  enemies.cheat: "true"
  enemies.cheat.level: noGoodRotten
  lives: "3"
  secret.code.allowed: "true"
  secret.code.passphrase: UUDDLRLRBABAS
kind: ConfigMap
metadata:
  name: game-env-config
  namespace: qae
  
# 从命令行直接输入 configmap 的内容
$ kubectl create configmap special-config --from-literal=special.how=very --from-literal=special.type=charm
```