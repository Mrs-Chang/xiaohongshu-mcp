# 小红书 MCP 服务 API 文档

## 概述

小红书 MCP (Model Context Protocol) 服务提供两套完整的 API 系统：

1. **MCP API 系统** - 使用 JSON-RPC 2.0 协议，专为 AI 客户端设计
2. **REST API 系统** - 标准 REST 风格 API，适合传统应用集成

## 服务架构

### 网络请求实现

是的，MCP 工具本质上就是网络请求。该服务通过以下方式实现：

- **服务端口**: `18060`
- **MCP 端点**: `http://localhost:18060/mcp`
- **REST API 端点**: `http://localhost:18060/api/v1/*`
- **协议支持**: JSON-RPC 2.0 (MCP) 和 HTTP REST
- **浏览器自动化**: 使用 Rod 框架控制浏览器与小红书网站交互

---

## 1. MCP API 系统

### 1.1 基本信息

- **端点**: `http://localhost:18060/mcp`
- **协议**: JSON-RPC 2.0 over HTTP (Streamable HTTP)
- **内容类型**: `application/json`
- **支持方法**: `POST`

### 1.2 请求格式

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "工具名称",
    "arguments": {
      "参数名": "参数值"
    }
  },
  "id": 1
}
```

### 1.3 响应格式

#### 成功响应
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "响应内容"
      }
    ],
    "isError": false
  },
  "id": 1
}
```

#### 错误响应
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "错误信息"
  },
  "id": 1
}
```

### 1.4 MCP 工具列表

#### 1.4.1 检查登录状态 - `check_login_status`

检查小红书登录状态。

**请求参数**: 无

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "check_login_status",
      "arguments": {}
    },
    "id": 1
  }'
```

**响应示例**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "登录状态检查成功: {\"is_logged_in\":true,\"username\":\"用户名\"}"
      }
    ]
  },
  "id": 1
}
```

#### 1.4.2 获取登录二维码 - `get_login_qrcode`

获取登录二维码，返回 Base64 图片和超时时间。

**请求参数**: 无

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "get_login_qrcode",
      "arguments": {}
    },
    "id": 1
  }'
```

**响应示例**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "请用小红书 App 在 2024-01-01 15:04:05 前扫码登录 👇"
      },
      {
        "type": "image",
        "mimeType": "image/png",
        "data": "iVBORw0KGgoAAAANSUhEUgAA..."
      }
    ]
  },
  "id": 1
}
```

#### 1.4.3 发布内容 - `publish_content`

发布图文内容到小红书。

**请求参数**:
- `title` (string, 必需): 内容标题，最多20个中文字或英文单词
- `content` (string, 必需): 正文内容，不包含以#开头的标签内容
- `images` (array, 必需): 图片路径列表，至少需要1张图片
- `tags` (array, 可选): 话题标签列表

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "publish_content",
      "arguments": {
        "title": "美食分享",
        "content": "今天做了一道美味的菜品",
        "images": ["/Users/user/image1.jpg", "https://example.com/image2.jpg"],
        "tags": ["美食", "生活", "分享"]
      }
    },
    "id": 1
  }'
```

#### 1.4.4 获取首页推荐 - `list_feeds`

获取小红书首页推荐内容列表。

**请求参数**: 无

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "list_feeds",
      "arguments": {}
    },
    "id": 1
  }'
```

#### 1.4.5 搜索内容 - `search_feeds`

搜索小红书内容，支持滚动加载更多结果。

**请求参数**:
- `keyword` (string, 必需): 搜索关键词
- `max_results` (string, 可选): 期望获取的最大结果数量
  - 正数: 指定具体数量（如 "30", "50"）
  - "-1": 尽可能多的结果（最多滚动10次）
  - 空或未提供: 使用默认值（约22个初始结果）

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "search_feeds",
      "arguments": {
        "keyword": "美食",
        "max_results": "50"
      }
    },
    "id": 1
  }'
```

#### 1.4.6 获取笔记详情 - `get_feed_detail`

获取小红书笔记详情，返回笔记内容、图片、作者信息、互动数据及评论列表。

**请求参数**:
- `feed_id` (string, 必需): 小红书笔记ID
- `xsec_token` (string, 必需): 访问令牌

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "get_feed_detail",
      "arguments": {
        "feed_id": "63f8d2a1000000001f01a2b3",
        "xsec_token": "XYZ123"
      }
    },
    "id": 1
  }'
```

#### 1.4.7 获取用户主页 - `user_profile`

获取小红书用户主页，返回用户基本信息，关注、粉丝、获赞量及其笔记内容。

**请求参数**:
- `user_id` (string, 必需): 小红书用户ID
- `xsec_token` (string, 必需): 访问令牌

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "user_profile",
      "arguments": {
        "user_id": "5f8d2a1000000001f01a2b3",
        "xsec_token": "XYZ123"
      }
    },
    "id": 1
  }'
```

#### 1.4.8 发表评论 - `post_comment_to_feed`

发表评论到小红书笔记。

**请求参数**:
- `feed_id` (string, 必需): 小红书笔记ID
- `xsec_token` (string, 必需): 访问令牌
- `content` (string, 必需): 评论内容

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "post_comment_to_feed",
      "arguments": {
        "feed_id": "63f8d2a1000000001f01a2b3",
        "xsec_token": "XYZ123",
        "content": "很棒的分享！"
      }
    },
    "id": 1
  }'
```

#### 1.4.9 下载无水印图片 - `download_images`

下载小红书笔记的无水印原图，支持多种格式。

**请求参数**:
- `feed_id` (string, 必需): 小红书笔记ID
- `xsec_token` (string, 必需): 访问令牌
- `format` (string, 可选): 图片格式，默认为 "png"
  - 支持: "png", "jpeg", "webp", "heic", "avif"
- `download_dir` (string, 可选): 下载目录路径，默认为 "downloads"

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "download_images",
      "arguments": {
        "feed_id": "63f8d2a1000000001f01a2b3",
        "xsec_token": "XYZ123",
        "format": "png",
        "download_dir": "downloads"
      }
    },
    "id": 1
  }'
```

#### 1.4.10 生成笔记链接 - `generate_url`

生成小红书笔记的完整链接，包含xsec_token参数。

**请求参数**:
- `feed_id` (string, 必需): 小红书笔记ID
- `xsec_token` (string, 可选): 访问令牌

**请求示例**:
```bash
curl -X POST http://localhost:18060/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "generate_url",
      "arguments": {
        "feed_id": "63f8d2a1000000001f01a2b3",
        "xsec_token": "XYZ123"
      }
    },
    "id": 1
  }'
```

---

## 2. REST API 系统

### 2.1 基本信息

- **基础 URL**: `http://localhost:18060/api/v1`
- **协议**: HTTP REST
- **内容类型**: `application/json`
- **认证**: 无需认证

### 2.2 通用响应格式

#### 成功响应
```json
{
  "success": true,
  "data": {},
  "message": "操作成功"
}
```

#### 错误响应
```json
{
  "error": "错误信息",
  "code": "ERROR_CODE",
  "details": "详细信息"
}
```

### 2.3 REST API 端点

#### 2.3.1 健康检查

**端点**: `GET /health`

**请求示例**:
```bash
curl -X GET http://localhost:18060/health
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "xiaohongshu-mcp",
    "account": "ai-report",
    "timestamp": "now"
  },
  "message": "服务正常"
}
```

#### 2.3.2 检查登录状态

**端点**: `GET /api/v1/login/status`

**请求示例**:
```bash
curl -X GET http://localhost:18060/api/v1/login/status
```

#### 2.3.3 获取登录二维码

**端点**: `GET /api/v1/login/qrcode`

**请求示例**:
```bash
curl -X GET http://localhost:18060/api/v1/login/qrcode
```

#### 2.3.4 发布内容

**端点**: `POST /api/v1/publish`

**请求体**:
```json
{
  "title": "内容标题",
  "content": "内容正文",
  "images": ["/path/to/image1.jpg"],
  "tags": ["标签1", "标签2"]
}
```

**请求示例**:
```bash
curl -X POST http://localhost:18060/api/v1/publish \
  -H "Content-Type: application/json" \
  -d '{
    "title": "美食分享",
    "content": "今天做了一道美味的菜品",
    "images": ["/Users/user/image1.jpg"],
    "tags": ["美食", "生活"]
  }'
```

#### 2.3.5 获取首页推荐

**端点**: `GET /api/v1/feeds/list`

**请求示例**:
```bash
curl -X GET http://localhost:18060/api/v1/feeds/list
```

#### 2.3.6 搜索内容

**端点**: `GET /api/v1/feeds/search`

**查询参数**:
- `keyword` (string, 必需): 搜索关键词
- `max_results` (string, 可选): 最大结果数量

**请求示例**:
```bash
curl -X GET "http://localhost:18060/api/v1/feeds/search?keyword=美食&max_results=50"
```

#### 2.3.7 获取笔记详情

**端点**: `POST /api/v1/feeds/detail`

**请求体**:
```json
{
  "feed_id": "笔记ID",
  "xsec_token": "访问令牌"
}
```

**请求示例**:
```bash
curl -X POST http://localhost:18060/api/v1/feeds/detail \
  -H "Content-Type: application/json" \
  -d '{
    "feed_id": "63f8d2a1000000001f01a2b3",
    "xsec_token": "XYZ123"
  }'
```

#### 2.3.8 获取用户主页

**端点**: `POST /api/v1/user/profile`

**请求体**:
```json
{
  "user_id": "用户ID",
  "xsec_token": "访问令牌"
}
```

**请求示例**:
```bash
curl -X POST http://localhost:18060/api/v1/user/profile \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "5f8d2a1000000001f01a2b3",
    "xsec_token": "XYZ123"
  }'
```

#### 2.3.9 发表评论

**端点**: `POST /api/v1/feeds/comment`

**请求体**:
```json
{
  "feed_id": "笔记ID",
  "xsec_token": "访问令牌",
  "content": "评论内容"
}
```

**请求示例**:
```bash
curl -X POST http://localhost:18060/api/v1/feeds/comment \
  -H "Content-Type: application/json" \
  -d '{
    "feed_id": "63f8d2a1000000001f01a2b3",
    "xsec_token": "XYZ123",
    "content": "很棒的分享！"
  }'
```

---

## 3. 数据类型定义

### 3.1 Feed 对象
```json
{
  "id": "笔记ID",
  "title": "笔记标题",
  "content": "笔记内容",
  "author": {
    "id": "作者ID",
    "name": "作者名称",
    "avatar": "头像URL"
  },
  "images": ["图片URL1", "图片URL2"],
  "tags": ["标签1", "标签2"],
  "stats": {
    "likes": 100,
    "comments": 50,
    "shares": 20
  },
  "created_at": "创建时间",
  "full_url": "完整链接"
}
```

### 3.2 用户信息对象
```json
{
  "userBasicInfo": {
    "id": "用户ID",
    "name": "用户名",
    "avatar": "头像URL",
    "description": "个人简介"
  },
  "interactions": [
    {
      "type": "followers",
      "count": 1000
    },
    {
      "type": "following",
      "count": 500
    },
    {
      "type": "likes",
      "count": 5000
    }
  ],
  "feeds": []
}
```

### 3.3 下载图片信息对象
```json
{
  "feed_id": "笔记ID",
  "title": "笔记标题",
  "total_images": 3,
  "downloaded_images": [
    {
      "index": 1,
      "original_url": "原始URL",
      "download_url": "下载URL",
      "local_path": "本地路径",
      "file_size": 1024000,
      "width": 1080,
      "height": 1920
    }
  ],
  "download_dir": "downloads",
  "format": "png"
}
```

---

## 4. 错误代码

### 4.1 JSON-RPC 错误代码
- `-32700`: 解析错误
- `-32600`: 无效请求
- `-32601`: 方法不存在
- `-32602`: 无效参数
- `-32603`: 内部错误

### 4.2 REST API 错误代码
- `STATUS_CHECK_FAILED`: 状态检查失败
- `INVALID_REQUEST`: 请求参数错误
- `PUBLISH_FAILED`: 发布失败
- `LIST_FEEDS_FAILED`: 获取列表失败
- `SEARCH_FEEDS_FAILED`: 搜索失败
- `GET_FEED_DETAIL_FAILED`: 获取详情失败
- `GET_USER_PROFILE_FAILED`: 获取用户信息失败
- `POST_COMMENT_FAILED`: 发表评论失败
- `MISSING_KEYWORD`: 缺少关键词参数

---

## 5. 注意事项

### 5.1 图片处理
- 支持 HTTP/HTTPS 链接（自动下载）
- 支持本地图片绝对路径（推荐）
- 发布时至少需要1张图片

### 5.2 参数限制
- 标题最多20个中文字或英文单词
- 内容正文不包含以#开头的标签内容
- 所有话题标签使用 `tags` 参数提供

### 5.3 登录状态
- 大部分操作需要先登录小红书
- 使用 `check_login_status` 检查登录状态
- 使用 `get_login_qrcode` 获取二维码登录

### 5.4 令牌获取
- `xsec_token` 从 Feed 列表或搜索结果中获取
- 用于访问需要认证的接口（详情、用户信息、评论等）

### 5.5 滚动搜索
- `max_results` 参数控制搜索结果数量
- 初始加载约22个结果
- 每次滚动加载10-20个额外结果
- 最多支持10次滚动

---

## 6. SDK 和工具

### 6.1 MCP 客户端支持
- Claude Desktop
- Cursor
- Cline
- Cherry Studio
- AnythingLLM
- 其他支持 HTTP MCP 的客户端

### 6.2 配置示例
```json
{
  "xiaohongshu-mcp": {
    "url": "http://localhost:18060/mcp",
    "type": "streamableHttp",
    "autoApprove": [],
    "disabled": false
  }
}
```

---

## 7. 技术实现

### 7.1 核心技术
- **Go 语言**: 后端服务实现
- **Gin 框架**: HTTP 服务器
- **Rod 框架**: 浏览器自动化
- **JSON-RPC 2.0**: MCP 协议实现
- **Streamable HTTP**: 支持流式响应

### 7.2 浏览器自动化
- 通过 Rod 控制 Chrome/Chromium 浏览器
- 模拟用户操作与小红书网站交互
- 支持无头模式和有头模式
- 自动处理 Cookie 和会话管理

### 7.3 图片处理
- 基于 XHS-Downloader 技术原理
- 支持无水印图片下载
- 多种格式转换支持
- 自动图片处理和优化

这个 API 文档涵盖了小红书 MCP 服务的所有网络请求接口，包括 MCP 工具和 REST API 两套完整的系统。每个接口都提供了详细的请求示例和响应格式说明。
