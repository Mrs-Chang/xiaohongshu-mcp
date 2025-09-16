#!/bin/bash

# 小红书 MCP 服务停止脚本

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

echo -e "${BLUE}正在停止小红书 MCP 服务...${NC}"

# 检查PID文件是否存在
if [ -f "$PID_FILE" ]; then
    pid=$(cat "$PID_FILE")
    echo -e "${YELLOW}找到PID文件，进程ID: $pid${NC}"
    
    # 检查进程是否还在运行
    if ps -p "$pid" > /dev/null 2>&1; then
        echo -e "${YELLOW}正在停止进程 $pid...${NC}"
        
        # 尝试优雅停止
        kill "$pid" 2>/dev/null || true
        sleep 3
        
        # 如果进程仍在运行，强制停止
        if ps -p "$pid" > /dev/null 2>&1; then
            echo -e "${YELLOW}进程仍在运行，强制停止...${NC}"
            kill -9 "$pid" 2>/dev/null || true
            sleep 1
        fi
        
        # 验证进程是否已停止
        if ps -p "$pid" > /dev/null 2>&1; then
            echo -e "${RED}❌ 无法停止进程 $pid${NC}"
        else
            echo -e "${GREEN}✅ 进程 $pid 已停止${NC}"
        fi
    else
        echo -e "${YELLOW}进程 $pid 已经不存在${NC}"
    fi
    
    # 删除PID文件
    rm -f "$PID_FILE"
else
    echo -e "${YELLOW}未找到PID文件${NC}"
fi

# 通过端口查找并停止服务
echo -e "${YELLOW}检查端口18060是否被占用...${NC}"
port_pids=$(lsof -ti:18060 2>/dev/null || true)
if [ -n "$port_pids" ]; then
    echo -e "${YELLOW}发现端口18060被以下进程占用: $port_pids${NC}"
    for port_pid in $port_pids; do
        echo -e "${YELLOW}停止进程 $port_pid...${NC}"
        kill -9 "$port_pid" 2>/dev/null || true
    done
    sleep 1
    
    # 再次检查端口
    port_pids=$(lsof -ti:18060 2>/dev/null || true)
    if [ -z "$port_pids" ]; then
        echo -e "${GREEN}✅ 端口18060已释放${NC}"
    else
        echo -e "${RED}❌ 端口18060仍被占用${NC}"
    fi
else
    echo -e "${GREEN}✅ 端口18060未被占用${NC}"
fi

# 通过进程名查找并停止相关进程
echo -e "${YELLOW}检查相关进程...${NC}"
pkill -f "xiaohongshu-mcp" 2>/dev/null && echo -e "${GREEN}✅ 已停止xiaohongshu-mcp相关进程${NC}" || true
pkill -f "go run.*main.go" 2>/dev/null && echo -e "${GREEN}✅ 已停止go run相关进程${NC}" || true

# 清理临时文件
echo -e "${YELLOW}清理临时文件...${NC}"
rm -f "$PID_FILE"

# 检查Chromium进程（可选）
chromium_pids=$(ps aux | grep -i chromium | grep -v grep | awk '{print $2}' || true)
if [ -n "$chromium_pids" ]; then
    echo -e "${YELLOW}发现Chromium进程，是否要停止？(y/N)${NC}"
    read -t 5 -r response || response="n"
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}停止Chromium进程...${NC}"
        echo "$chromium_pids" | xargs kill 2>/dev/null || true
        echo -e "${GREEN}✅ Chromium进程已停止${NC}"
    fi
fi

echo -e "${GREEN}🎉 小红书 MCP 服务已完全停止${NC}"

# 显示日志文件信息
if [ -f "$LOG_FILE" ]; then
    log_size=$(du -h "$LOG_FILE" | cut -f1)
    echo -e "${BLUE}日志文件: $LOG_FILE (大小: $log_size)${NC}"
    echo -e "${YELLOW}使用 'cat $LOG_FILE' 查看完整日志${NC}"
    echo -e "${YELLOW}使用 'rm $LOG_FILE' 删除日志文件${NC}"
fi
