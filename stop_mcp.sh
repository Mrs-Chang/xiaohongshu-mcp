#!/bin/bash

# 小红书MCP服务停止脚本
# 作者: xiaohongshu-mcp
# 用法: ./stop_mcp.sh

echo "🛑 停止小红书MCP服务..."

# 进入项目目录
cd "$(dirname "$0")"

# 检查是否有PID文件
if [ -f "mcp.pid" ]; then
    MCP_PID=$(cat mcp.pid)
    echo "📋 找到进程ID: $MCP_PID"
    
    # 尝试优雅停止
    if kill -TERM $MCP_PID 2>/dev/null; then
        echo "⏳ 正在优雅停止服务..."
        sleep 3
        
        # 检查进程是否还在运行
        if kill -0 $MCP_PID 2>/dev/null; then
            echo "⚠️  优雅停止失败，强制终止..."
            kill -9 $MCP_PID 2>/dev/null || true
        fi
    fi
    
    # 删除PID文件
    rm -f mcp.pid
else
    echo "⚠️  未找到PID文件，尝试查找并停止相关进程..."
fi

# 强制清理所有相关进程
echo "🧹 清理所有相关进程..."
pkill -9 -f "xiaohongshu-mcp" 2>/dev/null || true
pkill -9 -f "go run.*xiaohongshu-mcp" 2>/dev/null || true

# 清理端口占用
if lsof -i:18060 >/dev/null 2>&1; then
    echo "🔧 清理端口18060占用..."
    lsof -ti:18060 | xargs kill -9 2>/dev/null || true
fi

# 等待清理完成
sleep 2

# 验证清理结果
if lsof -i:18060 >/dev/null 2>&1; then
    echo "❌ 端口18060仍被占用，请手动检查"
    lsof -i:18060
    exit 1
elif pgrep -f "xiaohongshu-mcp" >/dev/null 2>&1; then
    echo "❌ 仍有相关进程在运行，请手动检查"
    pgrep -f "xiaohongshu-mcp"
    exit 1
else
    echo "✅ 小红书MCP服务已完全停止"
    echo "📄 日志文件保留在: $(pwd)/mcp.log"
    echo ""
    echo "💡 如需重新启动服务，请运行: ./start_mcp.sh"
fi
