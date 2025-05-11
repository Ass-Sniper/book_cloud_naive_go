
# **docker运行阶段app应用程序health健康检查失败问题**

```plaintext
kay@kay-vm:chap5$ docker compose down && docker compose up -d --build
Compose can now delegate builds to bake for better performance.
 To do so, set COMPOSE_BAKE=true.
[+] Building 22.5s (2/2) FINISHED                                                                                                                                      docker:default
 => [app internal] load build definition from Dockerfile                                                                                                                         0.0s
 => => transferring dockerfile: 1.17kB                                                                                                                                           0.0s
 => ERROR [app internal] load metadata for docker.io/library/ubuntu:20.04                                                                                                       22.4s
------
 > [app internal] load metadata for docker.io/library/ubuntu:20.04:
------
failed to solve: ubuntu:20.04: failed to resolve source metadata for docker.io/library/ubuntu:20.04: pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed
kay@kay-vm:chap5$
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker images
REPOSITORY   TAG       IMAGE ID   CREATED   SIZE
kay@kay-vm:chap5$
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04
20.04: Pulling from ddn-k8s/docker.io/ubuntu
560c024910be: Pull complete
Digest: sha256:38a0e8a00a21682240c31e48a3327dd7045dae42d300ff8e31e675660ac8dcbe
Status: Downloaded newer image for swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04
swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu:20.04  docker.io/ubuntu:20.04
kay@kay-vm:chap5$
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker images
REPOSITORY                                                  TAG       IMAGE ID       CREATED         SIZE
ubuntu                                                      20.04     5f5250218d28   11 months ago   72.8MB
swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/ubuntu   20.04     5f5250218d28   11 months ago   72.8MB
kay@kay-vm:chap5$
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker compose down && docker compose up -d --build
Compose can now delegate builds to bake for better performance.
 To do so, set COMPOSE_BAKE=true.
[+] Building 91.9s (27/27) FINISHED                                                                                                                                    docker:default
 => [app internal] load build definition from Dockerfile                                                                                                                         0.0s
 => => transferring dockerfile: 1.17kB                                                                                                                                           0.0s
 => [nginx internal] load metadata for docker.io/library/ubuntu:20.04                                                                                                            0.0s
 => [app internal] load .dockerignore                                                                                                                                            0.0s
 => => transferring context: 2B                                                                                                                                                  0.0s
 => [app internal] load build context                                                                                                                                            0.6s
 => => transferring context: 31.09MB                                                                                                                                             0.6s
 => CACHED [nginx 1/4] FROM docker.io/library/ubuntu:20.04                                                                                                                       0.0s
 => [app stage-1 2/6] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y ca-certificates tzd  9.7s
 => [app builder 2/8] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y wget curl git ca-c  31.3s
 => [app stage-1 3/6] WORKDIR /app                                                                                                                                               0.0s
 => [app builder 3/8] RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz &&     tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz                                           6.7s
 => [app builder 4/8] WORKDIR /app                                                                                                                                               0.0s
 => [app builder 5/8] COPY go.mod go.sum ./                                                                                                                                      0.1s
 => [app builder 6/8] RUN go mod download                                                                                                                                       16.0s
 => [app builder 7/8] COPY . ./                                                                                                                                                  0.1s
 => [app builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -o kvstore ./main.go                                                                                                17.2s
 => [app stage-1 4/6] COPY --from=builder /app/kvstore .                                                                                                                         0.1s
 => [app stage-1 5/6] COPY --from=builder /app/public ./public                                                                                                                   0.0s
 => [app stage-1 6/6] COPY --from=builder /app/users.txt ./users.txt                                                                                                             0.1s
 => [app] exporting to image                                                                                                                                                     0.4s
 => => exporting layers                                                                                                                                                          0.4s
 => => writing image sha256:eb930a7846e8f3f682a5748751c33adcc0f8e2fe3a787aea91d643ee2b3b9466                                                                                     0.0s
 => => naming to docker.io/library/chap5-app                                                                                                                                     0.0s
 => [app] resolving provenance for metadata file                                                                                                                                 0.0s
 => [nginx internal] load build definition from Dockerfile.nginx                                                                                                                 0.0s
 => => transferring dockerfile: 503B                                                                                                                                             0.0s
 => [nginx internal] load .dockerignore                                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                                  0.0s
 => [nginx internal] load build context                                                                                                                                          0.0s
 => => transferring context: 465B                                                                                                                                                0.0s
 => [nginx 2/4] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y nginx curl &&     rm -rf  18.9s
 => [nginx 3/4] COPY ./nginx.conf /etc/nginx/nginx.conf                                                                                                                          0.0s
 => [nginx 4/4] COPY ./public /usr/share/nginx/html                                                                                                                              0.0s
 => [nginx] exporting to image                                                                                                                                                   0.6s
 => => exporting layers                                                                                                                                                          0.6s
 => => writing image sha256:518915f25b92a14dac953e29fa2f174a68d17a8974bafecc4f5d1e5faf726241                                                                                     0.0s
 => => naming to docker.io/library/chap5-nginx                                                                                                                                   0.0s
 => [nginx] resolving provenance for metadata file                                                                                                                               0.0s
[+] Running 5/5
 ✔ app                        Built                                                                                                                                              0.0s
 ✔ nginx                      Built                                                                                                                                              0.0s
 ✔ Network chap5_kvstore_net  Created                                                                                                                                            0.1s
 ✘ Container chap5-app-1      Error                                                                                                                                             30.9s
 ✔ Container chap5-nginx-1    Created                                                                                                                                            0.0s
dependency failed to start: container chap5-app-1 is unhealthy
kay@kay-vm:chap5$
```

# **原因分析**

## **1. 进入对应容器手动检查** 

```shell
kay@kay-vm:chap5$
kay@kay-vm:chap5$ docker exec -it chap5-app-1 /bin/sh
#
# curl --help
/bin/sh: 2: curl: not found
# curl http://localhost:8080/health
/bin/sh: 3: curl: not found
# wget
/bin/sh: 4: wget: not found
# apk
/bin/sh: 5: apk: not found
# exit
kay@kay-vm:chap5$
```

## **2. 得出原因**
chap5\docker-compose.yml中配置了app的健康检查命令是curl，但在chap5\Dockerfile中【运行阶段】的镜像中未安装curl工具，从而导致app健康检查失败

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
    depends_on:
      app:
        condition: service_healthy
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    healthcheck:
      # 检查业务服务是否正常（8080 端口）
      test: ["CMD", "curl", "-f", "http://0.0.0.0:8080/health"]
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

# **修改chap5\Dockerfile后重新构建**

```plaintext
# 运行阶段
FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# 使用阿里源
RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -y wget curl ca-certificates tzdata && \    # <--- 此处增加wget和curl
    rm -rf /var/lib/apt/lists/*
```

```plaintext
kay@kay-vm:chap5$ docker compose down && docker compose up -d --build
[+] Running 3/3
 ✔ Container chap5-nginx-1    Removed                                                                                                                                            0.0s
 ✔ Container chap5-app-1      Removed                                                                                                                                            0.3s
 ✔ Network chap5_kvstore_net  Removed                                                                                                                                            0.3s
Compose can now delegate builds to bake for better performance.
 To do so, set COMPOSE_BAKE=true.
[+] Building 22.2s (27/27) FINISHED                                                                                                                                    docker:default
 => [app internal] load build definition from Dockerfile                                                                                                                         0.0s
 => => transferring dockerfile: 1.18kB                                                                                                                                           0.0s
 => [nginx internal] load metadata for docker.io/library/ubuntu:20.04                                                                                                            0.0s
 => [app internal] load .dockerignore                                                                                                                                            0.0s
 => => transferring context: 2B                                                                                                                                                  0.0s
 => [app internal] load build context                                                                                                                                            0.5s
 => => transferring context: 31.08MB                                                                                                                                             0.5s
 => [nginx 1/4] FROM docker.io/library/ubuntu:20.04                                                                                                                              0.0s
 => [app stage-1 2/6] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y wget curl ca-certi  11.7s
 => CACHED [app builder 2/8] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y wget curl gi  0.0s
 => CACHED [app builder 3/8] RUN wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz &&     tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz                                    0.0s
 => CACHED [app builder 4/8] WORKDIR /app                                                                                                                                        0.0s
 => CACHED [app builder 5/8] COPY go.mod go.sum ./                                                                                                                               0.0s
 => CACHED [app builder 6/8] RUN go mod download                                                                                                                                 0.0s
 => [app builder 7/8] COPY . ./                                                                                                                                                  0.2s
 => [app builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -o kvstore ./main.go                                                                                                20.9s
 => [app stage-1 3/6] WORKDIR /app                                                                                                                                               0.0s
 => [app stage-1 4/6] COPY --from=builder /app/kvstore .                                                                                                                         0.0s
 => [app stage-1 5/6] COPY --from=builder /app/public ./public                                                                                                                   0.0s
 => [app stage-1 6/6] COPY --from=builder /app/users.txt ./users.txt                                                                                                             0.0s
 => [app] exporting to image                                                                                                                                                     0.3s
 => => exporting layers                                                                                                                                                          0.3s
 => => writing image sha256:720ea3100928d6f2f95f6eee1590524f0ff5d126c50993d8a2fafe1dd2bf20b9                                                                                     0.0s
 => => naming to docker.io/library/chap5-app                                                                                                                                     0.0s
 => [app] resolving provenance for metadata file                                                                                                                                 0.0s
 => [nginx internal] load build definition from Dockerfile.nginx                                                                                                                 0.0s
 => => transferring dockerfile: 503B                                                                                                                                             0.0s
 => [nginx internal] load .dockerignore                                                                                                                                          0.0s
 => => transferring context: 2B                                                                                                                                                  0.0s
 => [nginx internal] load build context                                                                                                                                          0.0s
 => => transferring context: 465B                                                                                                                                                0.0s
 => CACHED [nginx 2/4] RUN sed -i 's|http://.*.ubuntu.com|http://mirrors.aliyun.com|g' /etc/apt/sources.list &&     apt-get update &&     apt-get install -y nginx curl &&       0.0s
 => CACHED [nginx 3/4] COPY ./nginx.conf /etc/nginx/nginx.conf                                                                                                                   0.0s
 => CACHED [nginx 4/4] COPY ./public /usr/share/nginx/html                                                                                                                       0.0s
 => [nginx] exporting to image                                                                                                                                                   0.0s
 => => exporting layers                                                                                                                                                          0.0s
 => => writing image sha256:518915f25b92a14dac953e29fa2f174a68d17a8974bafecc4f5d1e5faf726241                                                                                     0.0s
 => => naming to docker.io/library/chap5-nginx                                                                                                                                   0.0s
 => [nginx] resolving provenance for metadata file                                                                                                                               0.0s
[+] Running 5/5
 ✔ app                        Built                                                                                                                                              0.0s
 ✔ nginx                      Built                                                                                                                                              0.0s
 ✔ Network chap5_kvstore_net  Created                                                                                                                                            0.0s
 ✔ Container chap5-app-1      Healthy                                                                                                                                           10.8s
 ✔ Container chap5-nginx-1    Started                                                                                                                                           11.0s
kay@kay-vm:chap5$
```