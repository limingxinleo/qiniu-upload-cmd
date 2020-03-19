# 七牛云上传脚本

## 执行代码

修改 root.go 中的配置

然后执行一下脚本

```
go run main.go push data test
```

## 打包

```
docker build -t uploader .
docker cp <container id>:/go/builder/app uploader
```