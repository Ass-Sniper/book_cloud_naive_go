# 构建阶段
FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

# 设置阿里云源
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl git ca-certificates build-essential tzdata

# 安装 Go
RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:$PATH"

WORKDIR /app
COPY go.mod go.sum ./ 
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o kvstore ./main.go

# 运行阶段
FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# 使用阿里源
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/kvstore .
COPY --from=builder /app/public ./public
COPY --from=builder /app/users.txt ./users.txt
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/kvstore.db ./kvstore.db

EXPOSE 8080 6060
CMD ["./kvstore"]
