

启动api使用

```
docker build -t rustscanapi .
docker run -itd -p 50500:50500 rustscanapi ./rustapi-linux-amd64
```



```
http://ip:port/scan

[{
"ip":"127.0.0.1"
},{
"ip":"127.0.0.1"
}]
```



也可直接进入容器使用rustscan

```
dokcer exec -it [容器id] bash
```
