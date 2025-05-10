#!/bin/bash

# 停止名为 kvstore 的进程
echo "Stopping kvstore..."

pid=$(ps aux | grep '[k]vstore' | awk '{print $2}')

if [ -z "$pid" ]; then
  echo "kvstore is not running."
else
  kill "$pid"
  echo "kvstore (PID $pid) stopped."
fi

