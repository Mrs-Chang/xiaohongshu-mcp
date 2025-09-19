package xiaohongshu

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
)

func TestSearchWithScroll(t *testing.T) {
	t.Skip("SKIP: 测试滚动搜索功能")

	b := browser.NewBrowser(false)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := NewSearchAction(page)

	// 测试不滚动的搜索
	t.Run("不滚动搜索", func(t *testing.T) {
		feeds, err := action.SearchWithScroll(context.Background(), "美食", 0, 0)
		require.NoError(t, err)
		require.NotEmpty(t, feeds, "feeds should not be empty")

		fmt.Printf("不滚动搜索结果: %d 个帖子\n", len(feeds))
	})

	// 测试滚动3次的搜索
	t.Run("滚动3次搜索", func(t *testing.T) {
		feeds, err := action.SearchWithScroll(context.Background(), "美食", 3, 2*time.Second)
		require.NoError(t, err)
		require.NotEmpty(t, feeds, "feeds should not be empty")

		fmt.Printf("滚动3次搜索结果: %d 个帖子\n", len(feeds))
	})

	// 测试滚动5次的搜索
	t.Run("滚动5次搜索", func(t *testing.T) {
		feeds, err := action.SearchWithScroll(context.Background(), "美食", 5, 2*time.Second)
		require.NoError(t, err)
		require.NotEmpty(t, feeds, "feeds should not be empty")

		fmt.Printf("滚动5次搜索结果: %d 个帖子\n", len(feeds))

		// 打印前5个帖子的信息
		for i, feed := range feeds {
			if i >= 5 {
				break
			}
			fmt.Printf("帖子 %d: %s - %s\n", i+1, feed.NoteCard.DisplayTitle, feed.NoteCard.User.Nickname)
		}
	})
}

func TestScrollPerformance(t *testing.T) {
	t.Skip("SKIP: 测试滚动性能")

	b := browser.NewBrowser(false)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	action := NewSearchAction(page)

	// 比较不同滚动次数的效果
	testCases := []struct {
		name        string
		scrollCount int
		interval    time.Duration
	}{
		{"基础搜索", 0, 0},
		{"滚动1次", 1, 2 * time.Second},
		{"滚动3次", 3, 2 * time.Second},
		{"滚动5次", 5, 2 * time.Second},
		{"快速滚动5次", 5, 1 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			feeds, err := action.SearchWithScroll(context.Background(), "旅行", tc.scrollCount, tc.interval)
			duration := time.Since(start)

			require.NoError(t, err)
			require.NotEmpty(t, feeds)

			fmt.Printf("%s: %d 个帖子, 耗时: %v\n", tc.name, len(feeds), duration)
		})
	}
}
