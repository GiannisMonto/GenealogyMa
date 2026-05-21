#!/bin/bash
# ===================================
# 族谱数字化管理平台 - 启动脚本
# ===================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务端口
BACKEND_PORT=8080
ADMIN_PORT=3000
VISUALIZATION_PORT=5174

# 日志目录
LOG_DIR="./logs"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 创建日志目录
mkdir -p "$LOG_DIR"

# 打印标题
print_title() {
    echo -e "${BLUE}"
    echo "========================================"
    echo "   族谱数字化管理平台 - 启动脚本"
    echo "========================================"
    echo -e "${NC}"
}

# 检查端口是否可用
check_port() {
    local port=$1
    local name=$2
    if lsof -i :$port > /dev/null 2>&1; then
        echo -e "${YELLOW}[警告]${NC} 端口 $port ($name) 已被占用"
        return 1
    else
        echo -e "${GREEN}[OK]${NC} 端口 $port ($name) 可用"
        return 0
    fi
}

# 检查端口占用（Windows Git Bash / MSYS）
check_port_windows() {
    local port=$1
    local name=$2
    if netstat -ano 2>/dev/null | grep -q ":$port " || ss -ano 2>/dev/null | grep -q ":$port "; then
        echo -e "${YELLOW}[警告]${NC} 端口 $port ($name) 已被占用"
        return 1
    else
        echo -e "${GREEN}[OK]${NC} 端口 $port ($name) 可用"
        return 0
    fi
}

# 检查端口是否可用（跨平台）
check_port_crossplatform() {
    local port=$1
    local name=$2

    # 首先尝试 lsof
    if command -v lsof > /dev/null 2>&1; then
        if lsof -i :$port > /dev/null 2>&1; then
            echo -e "${YELLOW}[警告]${NC} 端口 $port ($name) 已被占用"
            return 1
        fi
    # 尝试 netstat (Windows)
    elif command -v netstat > /dev/null 2>&1; then
        if netstat -ano 2>/dev/null | grep -q ":$port "; then
            echo -e "${YELLOW}[警告]${NC} 端口 $port ($name) 已被占用"
            return 1
        fi
    fi

    echo -e "${GREEN}[OK]${NC} 端口 $port ($name) 可用"
    return 0
}

# 启动后端服务
start_backend() {
    echo -e "\n${BLUE}[1/3]${NC} 启动后端服务 (Go)..."
    cd backend
    nohup go run cmd/api/main.go > "$LOG_DIR/backend_$TIMESTAMP.log" 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > "$LOG_DIR/backend.pid"
    cd ..
    echo -e "${GREEN}[启动]${NC} 后端服务 PID: $BACKEND_PID"
    echo -e "${GREEN}[日志]${NC} $LOG_DIR/backend_$TIMESTAMP.log"
}

# 启动管理后台
start_admin() {
    echo -e "\n${BLUE}[2/3]${NC} 启动管理后台 (React)..."
    cd frontend/admin
    nohup npm run dev > "$LOG_DIR/admin_$TIMESTAMP.log" 2>&1 &
    ADMIN_PID=$!
    echo $ADMIN_PID > "$LOG_DIR/admin.pid"
    cd ../..
    echo -e "${GREEN}[启动]${NC} 管理后台 PID: $ADMIN_PID"
    echo -e "${GREEN}[日志]${NC} $LOG_DIR/admin_$TIMESTAMP.log"
}

# 启动可视化前端
start_visualization() {
    echo -e "\n${BLUE}[3/3]${NC} 启动可视化前端 (Vue 3)..."
    cd frontend/visualization
    nohup npm run dev > "$LOG_DIR/visualization_$TIMESTAMP.log" 2>&1 &
    VIS_PID=$!
    echo $VIS_PID > "$LOG_DIR/visualization.pid"
    cd ../..
    echo -e "${GREEN}[启动]${NC} 可视化前端 PID: $VIS_PID"
    echo -e "${GREEN}[日志]${NC} $LOG_DIR/visualization_$TIMESTAMP.log"
}

# 等待服务启动
wait_for_services() {
    echo -e "\n${BLUE}等待服务启动...${NC}"

    # 等待后端 (最多30秒)
    echo -n "等待后端服务 (端口 $BACKEND_PORT)..."
    for i in {1..30}; do
        if curl -s http://localhost:$BACKEND_PORT/api/health > /dev/null 2>&1 || \
           curl -s http://localhost:$BACKEND_PORT/swagger/index.html > /dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC}"
            break
        fi
        sleep 1
        echo -n "."
    done

    # 等待前端 (最多20秒)
    echo -n "等待管理后台 (端口 $ADMIN_PORT)..."
    for i in {1..20}; do
        if curl -s http://localhost:$ADMIN_PORT > /dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC}"
            break
        fi
        sleep 1
        echo -n "."
    done

    echo -n "等待可视化前端 (端口 $VISUALIZATION_PORT)..."
    for i in {1..20}; do
        if curl -s http://localhost:$VISUALIZATION_PORT > /dev/null 2>&1; then
            echo -e " ${GREEN}OK${NC}"
            break
        fi
        sleep 1
        echo -n "."
    done
}

# 打印启动信息
print_info() {
    echo -e "\n${GREEN}========================================"
    echo "   服务已启动!"
    echo "========================================${NC}"
    echo ""
    echo -e "  ${BLUE}管理后台:${NC}  http://localhost:$ADMIN_PORT"
    echo -e "  ${BLUE}可视化:${NC}    http://localhost:$VISUALIZATION_PORT"
    echo -e "  ${BLUE}API文档:${NC}   http://localhost:$BACKEND_PORT/swagger/index.html"
    echo ""
    echo -e "  ${YELLOW}停止服务:${NC}  ./scripts/stop.sh"
    echo -e "  ${YELLOW}查看日志:${NC}  tail -f $LOG_DIR/*.log"
    echo ""
}

# 主函数
main() {
    print_title

    echo -e "${BLUE}检查环境...${NC}"
    check_port_crossplatform $BACKEND_PORT "后端 API"
    check_port_crossplatform $ADMIN_PORT "管理后台"
    check_port_crossplatform $VISUALIZATION_PORT "可视化前端"

    echo -e "\n${BLUE}开始启动服务...${NC}"

    start_backend
    sleep 2
    start_admin
    sleep 1
    start_visualization

    wait_for_services
    print_info
}

# 清理函数
cleanup() {
    echo -e "\n${YELLOW}正在停止所有服务...${NC}"
    [ -f "$LOG_DIR/backend.pid" ] && kill $(cat "$LOG_DIR/backend.pid") 2>/dev/null && echo "后端已停止"
    [ -f "$LOG_DIR/admin.pid" ] && kill $(cat "$LOG_DIR/admin.pid") 2>/dev/null && echo "管理后台已停止"
    [ -f "$LOG_DIR/visualization.pid" ] && kill $(cat "$LOG_DIR/visualization.pid") 2>/dev/null && echo "可视化前端已停止"
    echo -e "${GREEN}所有服务已停止${NC}"
}

# 根据参数执行
case "${1:-start}" in
    start)
        main
        ;;
    stop)
        cleanup
        ;;
    restart)
        cleanup
        sleep 2
        main
        ;;
    status)
        echo -e "${BLUE}检查服务状态...${NC}"
        [ -f "$LOG_DIR/backend.pid" ] && kill -0 $(cat "$LOG_DIR/backend.pid") 2>/dev/null && echo -e "${GREEN}[运行中]${NC} 后端服务" || echo -e "${RED}[已停止]${NC} 后端服务"
        [ -f "$LOG_DIR/admin.pid" ] && kill -0 $(cat "$LOG_DIR/admin.pid") 2>/dev/null && echo -e "${GREEN}[运行中]${NC} 管理后台" || echo -e "${RED}[已停止]${NC} 管理后台"
        [ -f "$LOG_DIR/visualization.pid" ] && kill -0 $(cat "$LOG_DIR/visualization.pid") 2>/dev/null && echo -e "${GREEN}[运行中]${NC} 可视化前端" || echo -e "${RED}[已停止]${NC} 可视化前端"
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac