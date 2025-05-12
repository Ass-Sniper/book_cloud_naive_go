#!/bin/bash

set -e

echo "🔧 Starting project restructure..."

# 新目录结构
mkdir -p chap5/{cmd/kvstore,internal/{auth,handler,logger,store},web/{public,templates},config/{nginx,certs},scripts,docker,build,data}

# 移动 Go 源码
mv chap5/main.go chap5/cmd/kvstore/
mv chap5/auth/*.go chap5/internal/auth/
mv chap5/handler/*.go chap5/internal/handler/
mv chap5/logger/*.go chap5/internal/logger/
mv chap5/store/*.go chap5/internal/store/

# 移动前端文件
mv chap5/public/* chap5/web/public/
mv chap5/templates/* chap5/web/templates/

# 移动配置文件和证书
mv chap5/nginx.conf chap5/config/nginx/
mv chap5/Dockerfile.nginx chap5/config/nginx/
mv chap5/certs/* chap5/config/certs/

# 移动脚本
mv chap5/*.sh chap5/scripts/

# 移动 docker 配置
mv chap5/docker-compose.yml chap5/docker/
mv chap5/Dockerfile chap5/docker/

# 移动构建产物与数据
mv chap5/bin/* chap5/build/
mv chap5/kvstore.db chap5/data/
mv chap5/users.txt chap5/data/
mv chap5/boltbrowser.linux64 chap5/build/

# 移动 go.mod 和 go.sum
mv chap5/go.mod chap5/
mv chap5/go.sum chap5/

echo "✅ Project restructure completed."

