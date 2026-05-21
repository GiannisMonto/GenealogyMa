#!/bin/bash
# ===================================
# 族谱数字化管理平台 - 停止脚本
# ===================================

LOG_DIR="./logs"

echo -e "\033[94m========================================\033[0m"
echo -e "\033[94m   族谱数字化管理平台 - 停止脚本\033[0m"
echo -e "\033[94m========================================\033[0m"
echo ""

stop_service() {
    local name=$1
    local pid_file="$LOG_DIR/$name.pid"

    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null
            echo -e "\033[92m[停止]\033[0m $name 服务 (PID: $pid)"
        else
            echo -e "\033[93m[跳过]\033[0m $name 服务已停止"
        fi
        rm -f "$pid_file"
    else
        # 尝试通过进程名查找并终止
        local pids=$(pgrep -f "$name" 2>/dev/null || true)
        if [ -n "$pids" ]; then
            echo -e "\033[93m[警告]\033[0m 找到残留的 $name 进程: $pids"
            kill $pids 2>/dev/null || true
        else
            echo -e "\033[93m[跳过]\033[0m $name 服务未运行"
        fi
    fi
}

echo -e "\033[93m正在停止所有服务...\033[0m"
echo ""

stop_service "backend"
stop_service "admin"
stop_service "visualization"

echo ""
echo -e "\033[92m所有服务已停止\033[0m"