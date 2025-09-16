# 小红书无水印图片下载工具 (MCP)

## 概述

基于 XHS-Downloader 技术原理实现的 MCP (Model Context Protocol) 工具，用于下载小红书笔记的无水印原图。

## 技术原理

### 核心技术
1. **Token提取**: 从小红书图片URL中提取真实的文件标识符
2. **CDN直链构造**: 使用小红书CDN域名构造直接下载链接
3. **格式转换**: 通过URL参数控制输出格式
4. **绕过水印**: 直接访问原始文件，避开前端水印添加逻辑

### URL转换示例
```
原始URL: http://sns-webpic-qc.xhscdn.com/202509161717/674d2e41f2b3b972dc733d6258f3e27f/1040g2sg311v6icvpis6g49m398un3rak29ba220!nd_dft_wlteh_webp_3

转换后: https://ci.xiaohongshu.com/202509161717/674d2e41f2b3b972dc733d6258f3e27f/1040g2sg311v6icvpis6g49m398un3rak29ba220?imageView2/format/png
```

## MCP工具使用

### 工具名称
`download_images`

### 参数说明
- `feed_id` (必需): 小红书笔记ID，从Feed列表获取
- `xsec_token` (必需): 访问令牌，从Feed列表的xsecToken字段获取
- `format` (可选): 图片格式，默认为"png"
  - `png`: 无损高质量格式（推荐）
  - `jpeg`: 有损压缩格式
  - `webp`: 现代压缩格式
  - `heic`: 苹果设备格式
  - `avif`: 新一代格式
- `download_dir` (可选): 下载目录路径，默认为"downloads"

### 使用示例

#### 1. 搜索内容
```bash
# 通过MCP调用搜索工具
search_feeds(keyword="搞怪表情包")
```

#### 2. 获取详情
```bash
# 获取第一个帖子的详情
get_feed_detail(
    feed_id="66287569000000001c009d1e",
    xsec_token="ABGs9fTC7_0aCW7WKT-eXQg2IVobeuPXOeSm2TvUuatUM="
)
```

#### 3. 下载无水印图片
```bash
# 下载PNG格式的无水印原图
download_images(
    feed_id="66287569000000001c009d1e",
    xsec_token="ABGs9fTC7_0aCW7WKT-eXQg2IVobeuPXOeSm2TvUuatUM=",
    format="png",
    download_dir="downloads"
)
```

### 响应格式
```json
{
  "feed_id": "66287569000000001c009d1e",
  "title": "神金表情包",
  "total_images": 12,
  "downloaded_images": [
    {
      "index": 1,
      "original_url": "http://sns-webpic-qc.xhscdn.com/...",
      "download_url": "https://ci.xiaohongshu.com/...",
      "local_path": "downloads/神金表情包_1.png",
      "file_size": 245760,
      "width": 960,
      "height": 960
    }
  ],
  "download_dir": "downloads",
  "format": "png"
}
```

## 技术优势

### 1. 无水印下载
- 直接访问小红书CDN的原始文件
- 绕过前端水印添加逻辑
- 获取最高质量的原图

### 2. 格式灵活性
- 支持PNG、JPEG、WEBP、HEIC、AVIF等多种格式
- 通过CDN实时转换，无需本地处理
- 可根据需求选择最适合的格式

### 3. 高效稳定
- 基于官方数据结构，稳定性高
- 直接CDN访问，下载速度快
- 完善的错误处理和重试机制

### 4. 批量下载
- 自动下载笔记中的所有图片
- 智能文件命名（标题+序号）
- 并发下载提升效率

## 文件组织

### 下载目录结构
```
downloads/
├── 神金表情包_1.png
├── 神金表情包_2.png
├── 神金表情包_3.png
└── ...
```

### 文件命名规则
- 格式: `{笔记标题}_{序号}.{格式}`
- 自动清理不安全字符
- 长度限制防止文件系统错误
- 重复处理避免冲突

## 错误处理

### 常见错误及解决方案

1. **Feed不存在**: 检查feed_id和xsec_token是否正确
2. **没有图片**: 该笔记可能是视频或文本内容
3. **下载失败**: 网络问题或CDN访问限制
4. **格式不支持**: 使用支持的格式参数

### 日志监控
- 详细的操作日志记录
- 错误信息和堆栈追踪
- 下载进度和结果统计

## 注意事项

1. **合法使用**: 仅用于个人学习和研究目的
2. **版权尊重**: 下载的图片仍受原作者版权保护
3. **频率控制**: 避免过于频繁的请求，防止被限制
4. **存储管理**: 及时清理不需要的下载文件

## 相关技术

- [XHS-Downloader](https://github.com/JoeanAmier/XHS-Downloader): 技术原理参考
- [Model Context Protocol](https://modelcontextprotocol.io/): MCP协议规范
- [小红书CDN](https://ci.xiaohongshu.com/): 图片存储服务

---

*本工具基于XHS-Downloader项目的技术原理实现，仅用于技术学习和研究目的。*
