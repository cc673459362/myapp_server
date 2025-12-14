#!/bin/bash
# myapp_service.sh - 服务管理脚本

set -e  # 遇到错误退出

# 配置
APP_NAME="myapp_server"
APP_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"  # 脚本所在目录
BIN_DIR="$APP_DIR/bin"
BIN_NAME="$APP_NAME"
BIN_PATH="$BIN_DIR/$BIN_NAME"

# 日志和PID文件
LOG_DIR="$APP_DIR/logs"
PID_DIR="$APP_DIR/run"
PID_FILE="$PID_DIR/$APP_NAME.pid"
LOG_FILE="$LOG_DIR/$APP_NAME.log"
ERROR_LOG_FILE="$LOG_DIR/$APP_NAME.error.log"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 创建必要目录
create_dirs() {
    mkdir -p "$BIN_DIR" "$LOG_DIR" "$PID_DIR"
}

# 检查依赖
check_dependencies() {
    if [ ! -f "$BIN_PATH" ]; then
        echo -e "${RED}❌ 错误: 可执行文件不存在: $BIN_PATH${NC}"
        echo "请先运行构建脚本: ./build.sh"
        exit 1
    fi
    
    if [ ! -x "$BIN_PATH" ]; then
        chmod +x "$BIN_PATH"
    fi
}

# 检查进程是否运行
is_process_running() {
    local pid="$1"
    if [ -z "$pid" ]; then
        return 1
    fi
    
    if ps -p "$pid" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# 获取进程状态
get_process_status() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE" 2>/dev/null)
        if is_process_running "$pid"; then
            echo "running"
        else
            echo "stale"  # PID文件存在但进程不存在
        fi
    else
        echo "stopped"
    fi
}

# 显示状态
show_status() {
    local status=$(get_process_status)
    case "$status" in
        "running")
            local pid=$(cat "$PID_FILE")
            local uptime=$(ps -o etime= -p "$pid" 2>/dev/null | xargs || echo "unknown")
            echo -e "${GREEN}● $APP_NAME 运行中${NC}"
            echo "  PID: $pid"
            echo "  运行时间: $uptime"
            echo "  日志文件: $LOG_FILE"
            echo "  错误日志: $ERROR_LOG_FILE"
            return 0
            ;;
        "stale")
            echo -e "${YELLOW}⚠  $APP_NAME PID文件存在但进程未运行${NC}"
            echo "  PID文件: $PID_FILE"
            echo "  建议执行: $0 cleanup"
            return 2
            ;;
        "stopped")
            echo -e "${RED}○ $APP_NAME 未运行${NC}"
            return 3
            ;;
    esac
}

# 启动服务
start_service() {
    echo -e "${BLUE}▶ 启动 $APP_NAME 服务...${NC}"
    
    local status=$(get_process_status)
    if [ "$status" = "running" ]; then
        echo -e "${YELLOW}⚠  服务已在运行中${NC}"
        show_status
        return 0
    fi
    
    # 清理旧的PID文件
    if [ "$status" = "stale" ]; then
        echo -e "${YELLOW}⚠  清理旧的PID文件${NC}"
        rm -f "$PID_FILE"
    fi
    
    # 检查可执行文件
    check_dependencies
    
    # 检查是否已经有进程在运行（通过端口或其他方式）
    if lsof -ti:8080 >/dev/null 2>&1; then
        echo -e "${RED}❌ 错误: 端口 8080 已被占用${NC}"
        return 1
    fi
    
    # 切换到可执行文件目录
    cd "$BIN_DIR"
    
    # 记录启动时间
    echo "=== 服务启动于 $(date) ===" >> "$LOG_FILE"
    
    # 启动服务
    echo "启动命令: $BIN_PATH"
    echo "输出日志: $LOG_FILE"
    echo "错误日志: $ERROR_LOG_FILE"
    
    nohup "$BIN_PATH" >> "$LOG_FILE" 2>> "$ERROR_LOG_FILE" &
    local pid=$!
    
    # 保存PID
    echo $pid > "$PID_FILE"
    
    # 等待服务启动
    echo -n "等待服务启动"
    for i in {1..10}; do
        if is_process_running "$pid"; then
            # 检查服务是否就绪（假设HTTP端口8080）
            if curl -s http://localhost:8080/health >/dev/null 2>&1; then
                echo -e "\n${GREEN}✅ 服务启动成功！${NC}"
                show_status
                return 0
            fi
        fi
        echo -n "."
        sleep 1
    done
    
    echo -e "\n${YELLOW}⚠  服务已启动但健康检查未通过${NC}"
    show_status
    return 0
}

# 停止服务
stop_service() {
    echo -e "${BLUE}■ 停止 $APP_NAME 服务...${NC}"
    
    local status=$(get_process_status)
    if [ "$status" = "stopped" ]; then
        echo -e "${YELLOW}⚠  服务未运行${NC}"
        return 0
    fi
    
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        
        if is_process_running "$pid"; then
            echo "正在停止进程 $pid ..."
            
            # 先尝试优雅停止
            kill -TERM "$pid" 2>/dev/null
            
            # 等待最多10秒
            for i in {1..10}; do
                if ! is_process_running "$pid"; then
                    break
                fi
                echo -n "."
                sleep 1
            done
            
            # 如果还在运行，强制杀死
            if is_process_running "$pid"; then
                echo -e "\n${YELLOW}进程仍在运行，强制停止...${NC}"
                kill -9 "$pid" 2>/dev/null
                sleep 1
            fi
        fi
        
        # 清理PID文件
        if [ -f "$PID_FILE" ]; then
            rm -f "$PID_FILE"
        fi
        
        echo -e "${GREEN}✅ 服务已停止${NC}"
    else
        echo -e "${YELLOW}⚠  未找到PID文件${NC}"
    fi
}

# 重启服务
restart_service() {
    echo -e "${BLUE}🔄 重启 $APP_NAME 服务...${NC}"
    
    # 先停止
    if stop_service; then
        echo "等待2秒..."
        sleep 2
        # 再启动
        start_service
    else
        echo -e "${RED}❌ 停止服务失败${NC}"
        return 1
    fi
}

# 清理服务
cleanup_service() {
    echo -e "${BLUE}🧹 清理 $APP_NAME 服务...${NC}"
    
    # 1. 停止服务
    stop_service
    
    # 2. 清理PID文件
    if [ -f "$PID_FILE" ]; then
        rm -f "$PID_FILE"
        echo "已清理PID文件"
    fi
    
    # 3. 清理日志（可选，保留最近日志）
    if [ "$1" = "--all" ]; then
        echo "清理所有日志..."
        rm -f "$LOG_DIR"/*.log
    else
        # 只清理旧日志，保留最近3天
        find "$LOG_DIR" -name "*.log" -mtime +3 -delete
    fi
    
    # 4. 清理临时文件
    find "$APP_DIR/tmp" -type f -mtime +1 -delete 2>/dev/null || true
    
    echo -e "${GREEN}✅ 清理完成${NC}"
}

# 查看日志
view_logs() {
    local log_type="${1:-app}"
    
    case "$log_type" in
        "app"|"")
            echo -e "${BLUE}📄 查看应用日志:${NC}"
            if [ -f "$LOG_FILE" ]; then
                tail -f "$LOG_FILE"
            else
                echo -e "${YELLOW}日志文件不存在: $LOG_FILE${NC}"
            fi
            ;;
        "error")
            echo -e "${BLUE}📄 查看错误日志:${NC}"
            if [ -f "$ERROR_LOG_FILE" ]; then
                tail -f "$ERROR_LOG_FILE"
            else
                echo -e "${YELLOW}错误日志文件不存在: $ERROR_LOG_FILE${NC}"
            fi
            ;;
        "all")
            echo -e "${BLUE}📄 查看所有日志:${NC}"
            tail -f "$LOG_FILE" "$ERROR_LOG_FILE"
            ;;
        *)
            echo -e "${RED}❌ 未知日志类型: $log_type${NC}"
            echo "可用选项: app, error, all"
            return 1
            ;;
    esac
}

# 显示帮助
show_help() {
    echo -e "${BLUE}$APP_NAME 服务管理脚本${NC}"
    echo ""
    echo "用法: $0 {start|stop|restart|status|logs|cleanup|help}"
    echo ""
    echo "命令:"
    echo -e "  ${GREEN}start${NC}     启动服务"
    echo -e "  ${GREEN}stop${NC}      停止服务"
    echo -e "  ${GREEN}restart${NC}   重启服务"
    echo -e "  ${GREEN}status${NC}    查看服务状态"
    echo -e "  ${GREEN}logs${NC}      查看日志 (可加参数: app, error, all)"
    echo -e "  ${GREEN}cleanup${NC}   清理服务文件 (可加 --all 清理所有日志)"
    echo -e "  ${GREEN}help${NC}      显示此帮助"
    echo ""
    echo "示例:"
    echo "  $0 start          # 启动服务"
    echo "  $0 status         # 查看状态"
    echo "  $0 logs error     # 查看错误日志"
    echo "  $0 cleanup --all  # 彻底清理"
}

# 主函数
main() {
    # 创建必要目录
    create_dirs
    
    # 解析命令
    case "$1" in
        "start")
            start_service
            ;;
        "stop")
            stop_service
            ;;
        "restart")
            restart_service
            ;;
        "status")
            show_status
            ;;
        "logs")
            view_logs "$2"
            ;;
        "cleanup")
            cleanup_service "$2"
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        "")
            show_status
            ;;
        *)
            echo -e "${RED}❌ 未知命令: $1${NC}"
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"