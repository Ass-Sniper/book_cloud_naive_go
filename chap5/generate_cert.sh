#!/bin/bash

# 证书保存目录
CERT_DIR="./certs"
DOMAIN="192.168.16.248"  # 修改为你的 IP 或域名
DAYS=365
DHPARAM_FILE="$CERT_DIR/dhparam.pem"  # Diffie-Hellman 参数文件路径

# 创建目录
mkdir -p "$CERT_DIR"

# 生成证书和私钥
openssl req -x509 -nodes -newkey rsa:2048 \
  -days "$DAYS" \
  -keyout "$CERT_DIR/server.key" \
  -out "$CERT_DIR/server.crt" \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=Dev/OU=Local/CN=$DOMAIN"

# 生成 Diffie-Hellman 参数（可能需要一段时间）
openssl dhparam -out "$DHPARAM_FILE" 2048

# 提示信息
echo "✅ 自签名证书生成成功："
echo "  - 证书路径:     $CERT_DIR/server.crt"
echo "  - 私钥路径:     $CERT_DIR/server.key"
echo "  - 有效期:       $DAYS 天"
echo "  - CN（域名）:   $DOMAIN"
echo "  - Diffie-Hellman 参数: $DHPARAM_FILE"
