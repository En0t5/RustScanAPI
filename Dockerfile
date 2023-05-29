# 基础镜像
FROM debian:latest
# 维护者信息
MAINTAINER maintainer_name <test@test.com>

# 将应用程序代码复制到镜像中
COPY ./rustapi-linux-amd64 /root
COPY ./rustscan.deb /root
# 设置工作目录
WORKDIR /root

RUN apt-get update && \
    apt-get install -y nmap && \
    dpkg -i rustscan.deb && \
    chmod +x rustapi-linux-amd64

# 暴露端口
EXPOSE 50500

## 执行命令
#RUN nohup ./rustapi-linux-amd64
