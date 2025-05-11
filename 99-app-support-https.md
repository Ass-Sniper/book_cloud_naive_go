
# **1. 生成https证书**

```shell
kay@kay-vm:chap5$ ./generate_cert.sh
Generating a RSA private key
.+++++
.........................+++++
writing new private key to './certs/server.key'
-----
Generating DH parameters, 2048 bit long safe prime, generator 2
This is going to take a long time
...................................................................................................................................................................................................................+....................++*++*++*++*
✅ 自签名证书生成成功：
  - 证书路径:     ./certs/server.crt
  - 私钥路径:     ./certs/server.key
  - 有效期:       365 天
  - CN（域名）:   192.168.16.248
  - Diffie-Hellman 参数: ./certs/dhparam.pem
kay@kay-vm:chap5$
```


# **2. 修改配置文件**

## **2.1. 修改docker-compose.yml**

在docker-compose.yml中增加https证书的路径

```yaml
# docker-compose.yml

services:
  app:
    build: .
    environment:
      - TZ=Asia/Shanghai
    ports:
      - "8080:8080"  # 主业务接口
      - "6060:6060"  # pprof 接口
    networks:
      - kvstore_net
    healthcheck:
      # 检查业务服务是否正常（8080 端口）
      test: ["CMD", "curl", "-f", "http://0.0.0.0:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3

  nginx:
    build:
      context: .  # 使用自定义的 Dockerfile 构建 nginx 镜像
      dockerfile: Dockerfile.nginx  # 指定 Dockerfile 的名称
    ports:
      - "80:80"  # 将 nginx 对外暴露在 80 端口
      - "443:443"  # 将 nginx 对外暴露在 443 端口               # ✅ 加这行
    depends_on:
      app:
        condition: service_healthy
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./certs:/etc/nginx/certs  # 挂载证书                    # ✅ 加这行
    healthcheck:
      # 检查业务服务是否正常（80 端口）
      test: ["CMD", "curl", "-f", "http://0.0.0.0:80/health"]
      interval: 10s
      timeout: 5s
      retries: 3
    networks:
      - kvstore_net

networks:
  kvstore_net:
    driver: bridge
    attachable: true

```

## **2.2. 修改Dockerfile.nginx**

在Dockerfile.nginx中拷贝https证书的逻辑

```dockerfile
FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# 使用阿里云镜像源并安装 nginx 和 openssl
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y nginx curl openssl && \
    rm -rf /var/lib/apt/lists/*

# 复制自定义的 nginx 配置文件
COPY ./nginx.conf /etc/nginx/nginx.conf

# 复制前端文件
COPY ./public /usr/share/nginx/html

# 复制证书文件到容器内
COPY ./certs /etc/nginx/certs           # ✅ 加这行

# 暴露 80 和 443 端口（HTTP 和 HTTPS）
EXPOSE 80 443                           # ✅ 修改改行，增加443端口    

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]

```

## **2.2. 修改nginx.conf**

在nginx.conf中增加https的配置项

```plaintext
worker_processes auto;

events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;  # 容器服务名 + 端口
    }

    # ✅ 增加以下配置
    server {
        listen 80;
        # 强制所有 HTTP 请求重定向到 HTTPS
        return 301 https://$host$request_uri;   
    }

    # ✅ 修改配置
    server {
        listen 443 ssl http2;  # 启用 HTTP/2

        # SSL 配置
        ssl_certificate /etc/nginx/certs/server.crt;
        ssl_certificate_key /etc/nginx/certs/server.key;
        ssl_protocols TLSv1.2 TLSv1.3;  # 启用 TLSv1.2 和 TLSv1.3
        ssl_ciphers HIGH:!aNULL:!MD5;  # 配置加密套件
        ssl_prefer_server_ciphers on;  # 强制使用服务器端加密套件
        ssl_session_cache shared:SSL:10m;  # 启用 SSL 会话缓存
        ssl_session_timeout 1d;  # 设置 SSL 会话超时时间
        ssl_dhparam /etc/nginx/certs/dhparam.pem;  # 推荐使用 Diffie-Hellman 参数

        # 安全配置
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;

        # 健康检查路由
        location = /health {
            return 200 'OK';
            add_header Content-Type text/plain;
        }

        # 静态文件目录
        location /public/ {
            proxy_pass http://app/public/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_cache_valid 200 1d;
            add_header Cache-Control "public, max-age=86400";
        }

        # 默认代理业务接口
        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;  # 传递协议信息（http 或 https）
        }
    }
}

```


# **3. 重新构建docker镜像**

```shell
kay@kay-vm:chap5$ docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04
kay@kay-vm:chap5$ docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04  docker.io/ubuntu:20.04
kay@kay-vm:chap5$ sudo docker compose down && docker compose up -d --build
Compose can now delegate builds to bake for better performance.
 To do so, set COMPOSE_BAKE=true.
[+] Building 36.3s (30/30) FINISHED                                                                                                                                    docker:default
 => [app internal] load build definition from Dockerfile                                                                                                                         0.0s
 => => transferring dockerfile: 1.28kB                                                                                                                                           0.0s
 => [nginx internal] load metadata for docker.io/library/ubuntu:20.04                                                                                                            0.0s
 => [app internal] load .dockerignore                                                                                                                                            0.0s
 => => transferring context: 2B                                                                                                                                                  0.0s
 => [app internal] load build context                                                                                                                                            0.2s
 => => transferring context: 31.08MB                                                                                                                                             0.2s
 => CACHED [nginx 1/5] FROM docker.io/library/ubuntu:20.04                                                                                                                       0.0s
 => CACHED [app builder 2/8] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y wget curl gi  0.0s
 => CACHED [app builder 3/8] RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz &&     tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz                                    0.0s
 => CACHED [app builder 4/8] WORKDIR /app                                                                                                                                        0.0s
 => CACHED [app builder 5/8] COPY go.mod go.sum ./                                                                                                                               0.0s
 => CACHED [app builder 6/8] RUN go mod download                                                                                                                                 0.0s
 => [app builder 7/8] COPY . ./                                                                                                                                                  0.2s
 => [app builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -o kvstore ./main.go                                                                                                16.8s
 => CACHED [app stage-1 2/8] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y wget curl ca  0.0s
 => CACHED [app stage-1 3/8] WORKDIR /app                                                                                                                                        0.0s
 => CACHED [app stage-1 4/8] COPY --from=builder /app/kvstore .                                                                                                                  0.0s
 => CACHED [app stage-1 5/8] COPY --from=builder /app/public ./public                                                                                                            0.0s
 => CACHED [app stage-1 6/8] COPY --from=builder /app/users.txt ./users.txt                                                                                                      0.0s
 => CACHED [app stage-1 7/8] COPY --from=builder /app/templates ./templates                                                                                                      0.0s
 => CACHED [app stage-1 8/8] COPY --from=builder /app/kvstore.db ./kvstore.db                                                                                                    0.0s
 => [app] exporting to image                                                                                                                                                     0.0s
 => => exporting layers                                                                                                                                                          0.0s
 => => writing image sha256:363041f9055a63de0c7f339aba88db96692d8f220d2f26339c54dca4810993d9                                                                                     0.0s
 => => naming to docker.io/library/chap5-app                                                                                                                                     0.0s
 => [app] resolving provenance for metadata file                                                                                                                                 0.0s
 => [nginx internal] load build definition from Dockerfile.nginx                                                                                                                 0.0s
 => => transferring dockerfile: 716B                                                                                                                                             0.0s
 => [nginx internal] load .dockerignore                                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                                  0.0s
 => [nginx internal] load build context                                                                                                                                          0.0s
 => => transferring context: 599B                                                                                                                                                0.0s
 => [nginx 2/5] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y nginx curl openssl &&     18.1s
 => [nginx 3/5] COPY ./nginx.conf /etc/nginx/nginx.conf                                                                                                                          0.0s
 => [nginx 4/5] COPY ./public /usr/share/nginx/html                                                                                                                              0.0s
 => [nginx 5/5] COPY ./certs /etc/nginx/certs                                                                                                                                    0.0s
 => [nginx] exporting to image                                                                                                                                                   0.5s
 => => exporting layers                                                                                                                                                          0.5s
 => => writing image sha256:18a304e9d67ab797269c4a28c20029f07dfc44f0f0400b3843c9b4473a4a24e2                                                                                     0.0s
 => => naming to docker.io/library/chap5-nginx                                                                                                                                   0.0s
 => [nginx] resolving provenance for metadata file                                                                                                                               0.0s
[+] Running 5/5
 ✔ app                        Built                                                                                                                                              0.0s
 ✔ nginx                      Built                                                                                                                                              0.0s
 ✔ Network chap5_kvstore_net  Created                                                                                                                                            0.1s
 ✔ Container chap5-app-1      Healthy                                                                                                                                           10.8s
 ✔ Container chap5-nginx-1    Started                                                                                                                                           11.0s
kay@kay-vm:chap5$
```

# **4. 检查服务是否正常启动**

```shell
kay@kay-vm:chap5$ docker ps
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS                            PORTS                                                                                      NAMES
f849c9105278   chap5-nginx   "nginx -g 'daemon of…"   20 seconds ago   Up 8 seconds (health: starting)   0.0.0.0:80->80/tcp, [::]:80->80/tcp, 0.0.0.0:443->443/tcp, [::]:443->443/tcp               chap5-nginx-1
3582484dbcb6   chap5-app     "./kvstore"              20 seconds ago   Up 19 seconds (healthy)           0.0.0.0:6060->6060/tcp, [::]:6060->6060/tcp, 0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp   chap5-app-1
kay@kay-vm:chap5$
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker ps
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS                    PORTS                                                                                      NAMES
f849c9105278   chap5-nginx   "nginx -g 'daemon of…"   23 seconds ago   Up 11 seconds (healthy)   0.0.0.0:80->80/tcp, [::]:80->80/tcp, 0.0.0.0:443->443/tcp, [::]:443->443/tcp               chap5-nginx-1
3582484dbcb6   chap5-app     "./kvstore"              23 seconds ago   Up 22 seconds (healthy)   0.0.0.0:6060->6060/tcp, [::]:6060->6060/tcp, 0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp   chap5-app-1
kay@kay-vm:chap5$
```