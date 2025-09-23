#!/bin/bash

APP_DIR="/home/jiafengchen/go-projects/src/myapp_server"
# 进入项目目录（假设为项目根目录下的 cmd 子目录）
cd $APP_DIR/cmd || { echo "Error: 目录 cmd 不存在"; exit 1; }

# 执行构建（默认生成与目录同名的可执行文件）
go build -o myapp_server

# 检查构建结果
if [ $? -eq 0 ]; then
    echo "构建成功！可执行文件已生成在: $(pwd)"
else
    echo "构建失败，请检查错误日志"
    exit 1
fi
