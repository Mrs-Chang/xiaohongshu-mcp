package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

// ImageDownloader 图片下载器
type ImageDownloader struct {
	client *http.Client
}

// NewImageDownloader 创建图片下载器
func NewImageDownloader() *ImageDownloader {
	return &ImageDownloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// DownloadImages 下载图片列表
func (d *ImageDownloader) DownloadImages(ctx context.Context, imageList []xiaohongshu.DetailImageInfo, format, downloadDir, title string) ([]DownloadedImageInfo, error) {
	// 创建下载目录
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建下载目录失败: %w", err)
	}

	var downloadedImages []DownloadedImageInfo

	// 清理标题作为文件名前缀
	safeTitle := sanitizeFilename(title)

	for i, imageInfo := range imageList {
		// 提取图片token并构造无水印URL
		downloadURL, err := d.extractTokenAndBuildURL(imageInfo.URLDefault, format)
		if err != nil {
			logrus.Warnf("处理图片 %d 失败: %v", i+1, err)
			continue
		}

		// 构造本地文件路径
		filename := fmt.Sprintf("%s_%d.%s", safeTitle, i+1, format)
		localPath := filepath.Join(downloadDir, filename)

		// 下载图片
		fileSize, err := d.downloadFile(ctx, downloadURL, localPath)
		if err != nil {
			logrus.Warnf("下载图片 %d 失败: %v", i+1, err)
			continue
		}

		downloadedImages = append(downloadedImages, DownloadedImageInfo{
			Index:       i + 1,
			OriginalURL: imageInfo.URLDefault,
			DownloadURL: downloadURL,
			LocalPath:   localPath,
			FileSize:    fileSize,
			Width:       imageInfo.Width,
			Height:      imageInfo.Height,
		})

		logrus.Infof("成功下载图片 %d: %s", i+1, localPath)
	}

	if len(downloadedImages) == 0 {
		return nil, fmt.Errorf("没有成功下载任何图片")
	}

	return downloadedImages, nil
}

// extractTokenAndBuildURL 从原始URL中提取token并构造无水印CDN链接
// 基于XHS-Downloader的实现原理
func (d *ImageDownloader) extractTokenAndBuildURL(originalURL, format string) (string, error) {
	// 示例URL格式：
	// http://sns-webpic-qc.xhscdn.com/202509161717/674d2e41f2b3b972dc733d6258f3e27f/1040g2sg311v6icvpis6g49m398un3rak29ba220!nd_dft_wlteh_webp_3

	// 根据XHS-Downloader的实现，提取token的方式是：
	// 从URL中分割出路径部分，去除域名和查询参数
	// 然后去除!后面的后缀部分

	// 首先分割URL，提取路径部分（去除协议和域名）
	parts := strings.Split(originalURL, "/")
	if len(parts) < 6 {
		return "", fmt.Errorf("URL格式不正确: %s", originalURL)
	}

	// 重新组合路径部分（从第5个元素开始，跳过协议和域名）
	pathParts := parts[5:]
	fullPath := strings.Join(pathParts, "/")

	// 去除!后面的后缀
	if idx := strings.Index(fullPath, "!"); idx != -1 {
		fullPath = fullPath[:idx]
	}

	// 根据XHS-Downloader的实现，使用ci.xiaohongshu.com构造无水印链接
	// 格式：https://ci.xiaohongshu.com/{token}?imageView2/format/{format}
	downloadURL := fmt.Sprintf("https://ci.xiaohongshu.com/%s?imageView2/format/%s", fullPath, format)

	return downloadURL, nil
}

// downloadFile 下载文件到本地
func (d *ImageDownloader) downloadFile(ctx context.Context, url, localPath string) (int64, error) {
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	// 设置User-Agent模拟真实浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.xiaohongshu.com/")

	// 发送请求
	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 复制文件内容
	fileSize, err := io.Copy(file, resp.Body)
	if err != nil {
		return 0, err
	}

	return fileSize, nil
}

// sanitizeFilename 清理文件名，移除不安全字符
func sanitizeFilename(filename string) string {
	// 移除或替换不安全的文件名字符
	unsafe := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	safe := unsafe.ReplaceAllString(filename, "_")

	// 限制长度
	if len(safe) > 50 {
		safe = safe[:50]
	}

	// 移除首尾空格和点
	safe = strings.Trim(safe, " .")

	// 如果清理后为空，使用默认名称
	if safe == "" {
		safe = "image"
	}

	return safe
}
