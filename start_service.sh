#!/bin/bash

# 小红书 MCP 服务启动脚本
# 使用方法: ./start_service.sh [headless|gui]
# headless: 无头模式（默认）
# gui: 显示浏览器界面

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目目录
PROJECT_DIR="/Users/gaochang/xiaohongshu-mcp"
PID_FILE="$PROJECT_DIR/xiaohongshu-mcp.pid"
LOG_FILE="$PROJECT_DIR/xiaohongshu-mcp.log"

# 检查是否在项目目录
if [ ! -f "$PROJECT_DIR/main.go" ]; then
    echo -e "${RED}错误: 未找到项目文件，请检查项目路径${NC}"
    exit 1
fi

cd "$PROJECT_DIR"

# 检查服务是否已经在运行
check_service_running() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0  # 服务正在运行
        else
            rm -f "$PID_FILE"  # 删除无效的PID文件
            return 1  # 服务未运行
        fi
    fi
    return 1  # 服务未运行
}

# 停止已存在的服务
stop_existing_service() {
    echo -e "${YELLOW}检查是否有服务正在运行...${NC}"
    
    # 通过端口查找并停止服务
    port_pid=$(lsof -ti:18060 2>/dev/null || true)
    if [ -n "$port_pid" ]; then
        echo -e "${YELLOW}发现端口18060被占用，正在停止...${NC}"
        kill -9 "$port_pid" 2>/dev/null || true
        sleep 2
    fi
    
    # 通过进程名查找并停止服务
    pkill -f "xiaohongshu-mcp" 2>/dev/null || true
    pkill -f "go run.*main.go" 2>/dev/null || true
    
    # 清理PID文件
    rm -f "$PID_FILE"
    
    sleep 1
}

# 设置运行模式
MODE="headless"
if [ "$1" = "gui" ]; then
    MODE="gui"
    echo -e "${BLUE}启动模式: 显示浏览器界面${NC}"
elif [ "$1" = "headless" ] || [ -z "$1" ]; then
    MODE="headless"
    echo -e "${BLUE}启动模式: 无头模式（后台运行）${NC}"
else
    echo -e "${RED}错误: 无效的参数。使用 'headless' 或 'gui'${NC}"
    echo "使用方法: $0 [headless|gui]"
    exit 1
fi

# 停止已存在的服务
stop_existing_service

# 设置Go代理
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

echo -e "${BLUE}正在启动小红书 MCP 服务...${NC}"

# 根据模式启动服务
if [ "$MODE" = "headless" ]; then
    # 无头模式 - 后台运行
    nohup go run . -headless=true > "$LOG_FILE" 2>&1 &
    service_pid=$!
else
    # GUI模式 - 前台运行
    echo -e "${GREEN}服务将在前台运行，按 Ctrl+C 停止服务${NC}"
    go run . -headless=false
    exit 0
fi

# 保存PID
echo "$service_pid" > "$PID_FILE"

# 等待服务启动
echo -e "${YELLOW}等待服务启动...${NC}"
sleep 5

# 检查服务是否成功启动
if ps -p "$service_pid" > /dev/null 2>&1; then
    # 测试服务连接
    if curl -s -X POST http://localhost:18060/mcp \
       -H "Content-Type: application/json" \
       -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}' \
       > /dev/null 2>&1; then
        
        echo -e "${GREEN}✅ 小红书 MCP 服务启动成功！${NC}"
        echo -e "${GREEN}   服务地址: http://localhost:18060/mcp${NC}"
        echo -e "${GREEN}   PID: $service_pid${NC}"
        echo -e "${GREEN}   日志文件: $LOG_FILE${NC}"
        echo
        echo -e "${BLUE}可用的 MCP 工具:${NC}"
        echo "  • check_login_status - 检查登录状态"
        echo "  • publish_content - 发布图文内容"
        echo "  • list_feeds - 获取推荐列表"
        echo "  • search_feeds - 搜索内容"
        echo "  • get_feed_detail - 获取帖子详情"
        echo "  • post_comment_to_feed - 发表评论"
        echo
        echo -e "${YELLOW}使用 './stop_service.sh' 停止服务${NC}"
        echo -e "${YELLOW}使用 'tail -f $LOG_FILE' 查看实时日志${NC}"
    else
        echo -e "${RED}❌ 服务启动失败：无法连接到服务${NC}"
        kill "$service_pid" 2>/dev/null || true
        rm -f "$PID_FILE"
        exit 1
    fi
else
    echo -e "${RED}❌ 服务启动失败：进程异常退出${NC}"
    rm -f "$PID_FILE"
    echo -e "${YELLOW}查看日志获取详细错误信息: cat $LOG_FILE${NC}"
    exit 1
fi
