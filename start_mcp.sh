#!/bin/bash

# 小红书MCP服务启动脚本
# 作者: xiaohongshu-mcp
# 用法: ./start_mcp.sh

echo "🚀 启动小红书MCP服务..."

# 检查是否已有服务在运行
if lsof -i:18060 >/dev/null 2>&1; then
    echo "⚠️  端口18060已被占用，正在清理..."
    # 清理已存在的进程
    pkill -9 -f "xiaohongshu-mcp" 2>/dev/null || true
    pkill -9 -f "go run" 2>/dev/null || true
    lsof -ti:18060 | xargs kill -9 2>/dev/null || true
    sleep 2
fi

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go"
    exit 1
fi

# 进入项目目录
cd "$(dirname "$0")"

echo "📦 编译项目..."
if ! go build -o xiaohongshu-mcp .; then
    echo "❌ 编译失败"
    exit 1
fi

echo "🎯 启动服务 (非无头模式)..."
# 在后台启动服务
nohup ./xiaohongshu-mcp -headless=false > mcp.log 2>&1 &

# 获取进程ID
MCP_PID=$!
echo $MCP_PID > mcp.pid

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 3

# 检查服务是否成功启动
if curl -s -X POST http://localhost:18060/mcp \
   -H "Content-Type: application/json" \
   -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}' \
   | grep -q "xiaohongshu-mcp"; then
    echo "✅ 小红书MCP服务启动成功!"
    echo "📍 服务地址: http://localhost:18060/mcp"
    echo "📋 进程ID: $MCP_PID"
    echo "📄 日志文件: $(pwd)/mcp.log"
    echo ""
    echo "💡 使用方法:"
    echo "   - 在Cursor中已配置MCP，可直接使用"
    echo "   - 查看日志: tail -f mcp.log"
    echo "   - 停止服务: ./stop_mcp.sh"
else
    echo "❌ 服务启动失败，请查看日志文件: mcp.log"
    exit 1
fi
