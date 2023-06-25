

Docker 启动 API

```
docker build -t rustscanapi .
docker run -itd -p 50500:50500 rustscanapi ./rustapi-linux-amd64
```

API 接口

```python
# json 输入
POST http://ip:port/scan/json

[{
"ip":"127.0.0.1"
},{
"ip":"127.0.0.1"
}]

# text 输入
POST http://ip:port/scan/text

127.0.0.1
127.0.0.1

# 查看扫描结果
GET http://ip:port/show/result

# 下载扫描结果
GET http://ip:port/download/:filename

```


也可直接进入容器使用rustscan

```
dokcer exec -it [容器id] bash
```

