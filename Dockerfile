# 使用 golang 的镜像
FROM golang:1.18 as builderpodapi

# 设置工作目录
WORKDIR /app

# 复制代码到容器中
COPY . .

# 设置代理
ENV GOPROXY https://goproxy.cn,direct

# 编译应用程序
RUN go build -o main .

# 运行应用程序
CMD ["./main"]
