# ---------- 构建阶段 ----------
FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

# 设置阿里源并安装依赖
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl git ca-certificates build-essential tzdata

# 安装 Go
RUN wget https://golang.google.cn/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPROXY=https://goproxy.cn,direct

# 提前拷贝 go.mod 和 go.sum，用于缓存依赖层
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# 再拷贝剩余源码（能最大化利用缓存）
COPY . ./

# 构建可执行文件
RUN mkdir -p /build && CGO_ENABLED=0 GOOS=linux go build -o /build/kvstore ./cmd/kvstore

# ---------- 运行阶段 ----------
FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# 设置阿里源并清理缓存
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 复制构建好的二进制
COPY --from=builder /build/kvstore ./kvstore

# 复制前端资源
COPY --from=builder /app/web/public ./web/public
COPY --from=builder /app/web/templates ./web/templates

# 复制配置和初始化数据（可选）
COPY --from=builder /app/data ./data
COPY --from=builder /app/config/config.json ./config/config.json
COPY --from=builder /app/config/locales ./config/locales

# 开放端口
EXPOSE 8080 6060

# 启动命令
CMD ["./kvstore", "--config", "/app/config/config.json"]
    