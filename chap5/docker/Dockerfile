# 构建阶段
FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

# 设置阿里云源
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl git ca-certificates build-essential tzdata

# 安装 Go：使用国内镜像加速下载 Go
RUN wget https://golang.google.cn/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"


WORKDIR /app
# 拷贝 Go module 文件并下载依赖
COPY go.mod go.sum ./ 
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 拷贝所有源代码到构建目录，包括 web 目录

COPY . ./

# 构建可执行文件
RUN mkdir -p /build && CGO_ENABLED=0 GOOS=linux go build -o /build/kvstore ./cmd/kvstore

# 运行阶段
FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# 使用阿里源
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
# 只复制运行需要的文件
COPY --from=builder /build/kvstore ./kvstore
COPY --from=builder /app/web/public ./web/public
COPY --from=builder /app/web/templates ./web/templates

# 可选：不内置数据库（运行时挂载或首次初始化生成）
COPY --from=builder /app/data/users.txt ./data/users.txt
COPY --from=builder /app/data/kvstore.db ./data/kvstore.db

COPY --from=builder /app/config/config.json ./config/config.json

EXPOSE 8080 6060
CMD ["./kvstore", "--config", "/app/config/config.json"]
