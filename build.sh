#!/usr/bin/env bash
# build.sh - 智能构建脚本

set -e  # 遇到错误退出

# 自动检测项目根目录
find_project_root() {
    local dir="$PWD"
    while [ "$dir" != "/" ]; do
        if [ -f "$dir/go.mod" ]; then
            echo "$dir"
            return 0
        fi
        dir=$(dirname "$dir")
    done
    return 1
}

main() {
    echo "🔍 寻找 Go 项目根目录..."
    
    PROJECT_ROOT=$(find_project_root)
    if [ -z "$PROJECT_ROOT" ]; then
        echo "❌ 错误: 未找到 go.mod 文件，当前目录不是 Go 项目"
        exit 1
    fi
    
    echo "📁 项目根目录: $PROJECT_ROOT"
    cd "$PROJECT_ROOT"
    
    # 配置
    local APP_NAME="myapp_server"
    local CMD_PATH="./cmd"
    local BIN_DIR="./bin"
    local VERSION="1.0.0"
    
    # 清理
    echo "🧹 清理..."
    rm -rf "$BIN_DIR"
    mkdir -p "$BIN_DIR"
    
    # 构建
    echo "🔨 构建 $APP_NAME..."
    if go build -o "$BIN_DIR/$APP_NAME" "$CMD_PATH"; then
        echo "✅ 构建成功！"
        echo ""
        echo "📊 构建信息:"
        echo "   可执行文件: $BIN_DIR/$APP_NAME"
        echo "   文件大小: $(du -h "$BIN_DIR/$APP_NAME" | cut -f1)"
        echo ""
        echo "🚀 运行命令:"
        echo "   $BIN_DIR/$APP_NAME"
        echo "   或"
        echo "   ./$BIN_DIR/$APP_NAME"
    else
        echo "❌ 构建失败"
        exit 1
    fi
}

# 运行
main "$@"