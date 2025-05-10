#!/bin/bash

# 设置项目名称和输出目录
PROJECT_NAME="kvstore"
OUTPUT_DIR="./bin"

# 检查 go.mod 文件是否存在
if [ ! -f go.mod ]; then
  echo "go.mod 文件不存在，正在初始化 Go 模块..."
  # 初始化 Go 模块
  go mod init github.com/kay/$PROJECT_NAME
fi

# 清理并下载依赖
echo "正在清理并下载依赖..."
rm -rf ./bin/*
go mod tidy

go mod download

# 创建 bin 输出目录
mkdir -p $OUTPUT_DIR

# 构建项目
echo "正在构建项目..."

if [ -d .git ]; then
  BUILD_FLAGS=""
else 
  BUILD_FLAGS="-buildvcs=false"
fi

# 编译当前平台的二进制文件
go build -o $OUTPUT_DIR/$PROJECT_NAME $BUILD_FLAGS 

# 如果需要支持多平台编译，可以在这里添加不同平台的构建
# 以下为示例，可以在 Linux、macOS、Windows 下交叉编译

# 编译 Linux 版本
echo "正在编译 Linux 版本..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$PROJECT_NAME-linux $BUILD_FLAGS

# 打印编译完成的信息
echo "构建完成！可执行文件已存放在 $OUTPUT_DIR 目录下"
echo "Linux 版本: $OUTPUT_DIR/$PROJECT_NAME-linux"
echo "当前平台版本: $OUTPUT_DIR/$PROJECT_NAME"
