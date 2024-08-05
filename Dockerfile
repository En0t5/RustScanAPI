# 基础镜像
FROM hub.atomgit.com/amd64/debian:rc-buggy
# 维护者信息
MAINTAINER RustScanAPI <test@test.com>

# 添加项目路径和缓存路径
RUN mkdir /app && mkdir /app/cache

# 将应用程序代码复制到镜像中
COPY ./rustapi-linux-amd64 /app
COPY ./rustscan.deb /app
# 设置工作目录
WORKDIR /app

# 设置国内源
RUN cp /etc/apt/sources.list /etc/apt/sources.list.bak &  echo "deb http://mirrors.aliyun.com/debian/ stable main contrib non-free" > /etc/apt/sources.list && \
                                                          echo "deb-src http://mirrors.aliyun.com/debian/ stable main contrib non-free" >> /etc/apt/sources.list && \
                                                          echo "deb http://mirrors.aliyun.com/debian-security/ stable-security main contrib non-free" >> /etc/apt/sources.list && \
                                                          echo "deb-src http://mirrors.aliyun.com/debian-security/ stable-security main contrib non-free" >> /etc/apt/sources.list && \
                                                          echo "deb http://mirrors.aliyun.com/debian/ stable-updates main contrib non-free" >> /etc/apt/sources.list && \
                                                          echo "deb-src http://mirrors.aliyun.com/debian/ stable-updates main contrib non-free" >> /etc/apt/sources.list


RUN apt-get update && \
    apt-get install -y nmap && \
    dpkg -i rustscan.deb && \
    chmod +x rustapi-linux-amd64

# 暴露端口
EXPOSE 50500

## 执行命令
#RUN nohup ./rustapi-linux-amd64
