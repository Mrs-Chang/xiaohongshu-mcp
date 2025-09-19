package xiaohongshu

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
)

func TestSearch(t *testing.T) {

	t.Skip("SKIP: 测试发布")

	b := browser.NewBrowser(false)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := NewSearchAction(page)

	feeds, err := action.Search(context.Background(), "Kimi")
	require.NoError(t, err)
	require.NotEmpty(t, feeds, "feeds should not be empty")

	fmt.Printf("成功获取到 %d 个 Feed\n", len(feeds))

	for _, feed := range feeds {
		fmt.Printf("Feed ID: %s\n", feed.ID)
		fmt.Printf("Feed Title: %s\n", feed.NoteCard.DisplayTitle)
	}
}

// TestSearchEmoji 测试搜索表情包相关帖子
func TestSearchEmoji(t *testing.T) {
	// 不跳过这个测试，让它实际运行
	// t.Skip("SKIP: 测试表情包搜索")

	b := browser.NewBrowser(false)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := NewSearchAction(page)

	// 测试基础搜索
	t.Run("基础搜索表情包", func(t *testing.T) {
		feeds, err := action.Search(context.Background(), "表情包")
		require.NoError(t, err)
		require.NotEmpty(t, feeds, "feeds should not be empty")

		fmt.Printf("=== 基础搜索 '表情包' 结果 ===\n")
		fmt.Printf("共找到 %d 个帖子\n\n", len(feeds))

		// 显示前10个帖子的详细信息
		for i, feed := range feeds {
			if i >= 10 {
				break
			}
			fmt.Printf("--- 帖子 %d ---\n", i+1)
			fmt.Printf("ID: %s\n", feed.ID)
			fmt.Printf("标题: %s\n", feed.NoteCard.DisplayTitle)
			fmt.Printf("作者: %s\n", feed.NoteCard.User.Nickname)
			fmt.Printf("点赞数: %s\n", feed.NoteCard.InteractInfo.LikedCount)
			fmt.Printf("评论数: %s\n", feed.NoteCard.InteractInfo.CommentCount)
			fmt.Printf("收藏数: %s\n", feed.NoteCard.InteractInfo.CollectedCount)
			fmt.Printf("完整链接: %s\n", feed.GetFullURL())
			fmt.Printf("\n")
		}
	})

	// 测试滚动搜索
	t.Run("滚动搜索表情包", func(t *testing.T) {
		feeds, err := action.SearchWithScroll(context.Background(), "表情包", 3, 2*time.Second) // 2秒
		require.NoError(t, err)
		require.NotEmpty(t, feeds, "feeds should not be empty")

		fmt.Printf("=== 滚动搜索 '表情包' 结果 (滚动3次) ===\n")
		fmt.Printf("共找到 %d 个帖子\n\n", len(feeds))

		// 显示统计信息
		var videoCount, imageCount int
		for _, feed := range feeds {
			if feed.NoteCard.Type == "video" {
				videoCount++
			} else {
				imageCount++
			}
		}

		fmt.Printf("内容类型统计:\n")
		fmt.Printf("- 视频帖子: %d 个\n", videoCount)
		fmt.Printf("- 图片帖子: %d 个\n", imageCount)
		fmt.Printf("\n")

		// 显示点赞数最高的5个帖子
		fmt.Printf("=== 热门帖子 (按点赞数排序) ===\n")
		// 简单排序，找出点赞数最高的几个
		for i := 0; i < len(feeds) && i < 5; i++ {
			feed := feeds[i]
			fmt.Printf("%d. %s (点赞: %s)\n", i+1, feed.NoteCard.DisplayTitle, feed.NoteCard.InteractInfo.LikedCount)
		}
	})
}
