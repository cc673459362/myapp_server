#!/bin/bash
APP_NAME="myapp_server"
APP_DIR="/home/jiafengchen/go-projects"
BIN_PATH="$APP_DIR/bin/myapp_server"
CONFIG_PATH="$APP_DIR/config/app.env"
LOG_FILE="$APP_DIR/logs/myapp_server.log"
PID_FILE="$APP_DIR/PIDFILE"

# 确保存在日志目录
mkdir -p "$APP_DIR/logs"

start_service() {
    if [ -f "$PID_FILE" ]; then
        echo "服务已在运行中 (PID: $(cat $PID_FILE))"
        exit 1
    fi

    cd $APP_DIR/bin
    echo "启动 $APP_NAME ..."
    nohup "$BIN_PATH" >> "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    echo "服务已启动，PID: $(cat $PID_FILE)"
}

stop_service() {
    if [ ! -f "$PID_FILE" ]; then
        echo "服务未运行"
        exit 1
    fi

    PID=$(cat "$PID_FILE")
    echo "停止 $APP_NAME (PID: $PID)..."
    kill -SIGTERM "$PID"
    rm "$PID_FILE"
    echo "服务已停止"
}

restart_service() {
    stop_service
    sleep 2
    start_service
}

service_status() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null; then
            echo "● $APP_NAME 运行中 (PID: $PID)"
            return 0
        else
            echo "● $APP_NAME 进程存在但未运行"
            return 2
        fi
    else
        echo "○ $APP_NAME 未运行"
        return 3
    fi
}

case "$1" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    status)
        service_status
        ;;
    *)
        echo "使用方法: $0 {start|stop|restart|status}"
        exit 1
esac
