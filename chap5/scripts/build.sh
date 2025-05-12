#!/bin/bash
set -e

# 设置项目名称和输出目录
PROJECT_NAME="kvstore"
OUTPUT_DIR="./build"

# 检查 go.mod 文件是否存在
if [ ! -f go.mod ]; then
  echo "go.mod 文件不存在，正在初始化 Go 模块..."
  go mod init github.com/kay/$PROJECT_NAME
fi

# 清理并下载依赖
echo "正在清理并下载依赖..."
rm -rf $OUTPUT_DIR/$PROJECT_NAME
go mod tidy
go mod download

# 创建输出目录
mkdir -p $OUTPUT_DIR

# 项目主入口
MAIN_PKG="./cmd/kvstore"

# 是否禁用 VCS 构建信息
if [ -d .git ]; then
  BUILD_FLAGS=""
else
  BUILD_FLAGS="-buildvcs=false"
fi

# 构建当前平台的二进制文件
echo "正在构建当前平台版本..."
go build -o $OUTPUT_DIR/$PROJECT_NAME $BUILD_FLAGS $MAIN_PKG

# 打印构建完成信息
echo "构建完成！可执行文件路径：$OUTPUT_DIR/$PROJECT_NAME"
