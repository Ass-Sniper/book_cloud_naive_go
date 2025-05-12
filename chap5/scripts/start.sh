#!/bin/bash

# 设置项目名称和二进制文件路径
PROJECT_NAME="kvstore"
OUTPUT_DIR="./bin"
BINARY="$OUTPUT_DIR/$PROJECT_NAME"

# 检查二进制文件是否存在
if [ ! -f $BINARY ]; then
  echo "错误: 未找到可执行文件 $BINARY，请先运行 'build.sh' 进行构建。"
  exit 1
fi

# 启动服务
echo "正在启动 $PROJECT_NAME 服务..."
nohup $BINARY > $OUTPUT_DIR/$PROJECT_NAME.log 2>&1 &

# 输出启动日志
echo "服务已启动，日志输出在 $OUTPUT_DIR/$PROJECT_NAME.log"
echo "访问地址: http://localhost:8080"
