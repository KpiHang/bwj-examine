# 使用官方 Go 镜像作为构建环境
FROM golang:1.23 AS builder
# 设置工作目录
WORKDIR /app

# 复制 src 目录下的所有文件到工作目录
COPY ./src2 .

# 构建 Go 应用
RUN go build -o myapp

# 使用一个更小的基础镜像
FROM alpine:latest

# 安装 curl
RUN apk --no-cache add ca-certificates curl

# 设置工作目录
WORKDIR /root/

# 从构建环境中复制二进制文件到当前镜像
COPY --from=builder /app/myapp .

# 运行应用
CMD ["./myapp"]
