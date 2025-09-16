package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// handleGenerateURL 处理生成小红书完整链接
func (s *AppServer) handleGenerateURL(ctx context.Context, args map[string]interface{}) *MCPToolResult {
	logrus.Info("MCP: 生成小红书链接")

	// 解析参数
	feedID, ok := args["feed_id"].(string)
	if !ok || feedID == "" {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "生成链接失败: 缺少feed_id参数",
			}},
			IsError: true,
		}
	}

	xsecToken, _ := args["xsec_token"].(string) // xsec_token 是可选的

	logrus.Infof("MCP: 生成链接 - Feed ID: %s, XsecToken: %s", feedID, xsecToken)

	// 生成基础链接
	baseURL := "https://www.xiaohongshu.com/explore/" + feedID
	fullURL := baseURL
	
	// 如果有 xsec_token，添加到链接中
	if xsecToken != "" {
		fullURL = baseURL + "?xsec_token=" + xsecToken
	}

	// 返回结果
	result := map[string]string{
		"feed_id":    feedID,
		"base_url":   baseURL,
		"full_url":   fullURL,
		"xsec_token": xsecToken,
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("生成链接成功，但序列化失败: %v", err),
			}},
			IsError: true,
		}
	}

	return &MCPToolResult{
		Content: []MCPContent{{
			Type: "text",
			Text: string(jsonData),
		}},
	}
}
