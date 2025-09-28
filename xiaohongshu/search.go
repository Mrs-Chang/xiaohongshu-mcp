package xiaohongshu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-rod/rod"
)

type SearchResult struct {
	Search struct {
		Feeds FeedsValue `json:"feeds"`
	} `json:"search"`
}

type SearchAction struct {
	page *rod.Page
}

func NewSearchAction(page *rod.Page) *SearchAction {
	pp := page.Timeout(60 * time.Second)

	return &SearchAction{page: pp}
}

func (s *SearchAction) Search(ctx context.Context, keyword string) ([]Feed, error) {
	return s.SearchWithScroll(ctx, keyword, 0, 0)
}

// SearchWithScroll 带滚动功能的搜索，可以加载更多内容
func (s *SearchAction) SearchWithScroll(ctx context.Context, keyword string, scrollCount int, scrollInterval time.Duration) ([]Feed, error) {
	page := s.page.Context(ctx)

	searchURL := makeSearchURL(keyword)
	page.MustNavigate(searchURL)
	page.MustWaitStable()

	page.MustWait(`() => window.__INITIAL_STATE__ !== undefined`)

	// 如果需要滚动加载更多内容
	if scrollCount > 0 {
		if err := s.performScrolling(page, scrollCount, scrollInterval); err != nil {
			return nil, fmt.Errorf("scrolling failed: %w", err)
		}
	}

	// 获取 window.__INITIAL_STATE__ 并转换为 JSON 字符串
	result := page.MustEval(`() => {
			if (window.__INITIAL_STATE__) {
				return JSON.stringify(window.__INITIAL_STATE__);
			}
			return "";
		}`).String()

	if result == "" {
		return nil, fmt.Errorf("__INITIAL_STATE__ not found")
	}

	var searchResult SearchResult
	if err := json.Unmarshal([]byte(result), &searchResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal __INITIAL_STATE__: %w", err)
	}

	return searchResult.Search.Feeds.Value, nil
}

// performScrolling 执行连续滚动操作
func (s *SearchAction) performScrolling(page *rod.Page, scrollCount int, scrollInterval time.Duration) error {
	// 设置默认滚动间隔
	if scrollInterval == 0 {
		scrollInterval = time.Second
	}

	// 获取初始的帖子数量，用于监控是否有新内容加载
	initialCount := s.getCurrentFeedCount(page)
	previousCount := initialCount // 保存上一次的数量，用于比较

	fmt.Printf("开始滚动加载更多内容，初始帖子数量: %d\n", initialCount)

	for i := 0; i < scrollCount; i++ {
		fmt.Printf("第 %d/%d 次滚动...\n", i+1, scrollCount)

		// 执行滚动操作
		_, err := page.Eval(`() => {
			// 获取滚动前的高度
			const beforeHeight = document.body.scrollHeight;
			console.log('滚动前页面高度:', beforeHeight);
			
			// 滚动到页面底部
			window.scrollTo({
				top: document.body.scrollHeight,
				behavior: 'smooth'
			});
			
			// 记录滚动信息
			console.log('已滚动到位置:', window.pageYOffset);
			console.log('滚动后页面高度:', document.body.scrollHeight);
			
			return {
				scrollPosition: window.pageYOffset,
				beforeHeight: beforeHeight,
				afterHeight: document.body.scrollHeight
			};
		}`)

		if err != nil {
			return fmt.Errorf("scroll operation failed: %w", err)
		}

		// 等待内容加载 - 增加等待时间
		fmt.Printf("等待 %v 让内容加载...\n", scrollInterval)
		time.Sleep(scrollInterval)

		// 再等待一点时间确保懒加载完成
		time.Sleep(1 * time.Second)

		// 等待页面稳定
		page.MustWaitStable()

		// 检查是否有新内容加载
		currentCount := s.getCurrentFeedCount(page)
		fmt.Printf("滚动后帖子数量: %d (新增: %d)\n", currentCount, currentCount-initialCount)

		// 移除提前停止逻辑，让滚动执行完整的次数
		previousCount = currentCount // 更新上一次的数量
	}

	fmt.Printf("滚动完成，最终帖子数量: %d\n", previousCount)
	return nil
}

// getCurrentFeedCount 获取当前页面的帖子数量
func (s *SearchAction) getCurrentFeedCount(page *rod.Page) int {
	// 尝试多种方式获取帖子数量，以确保准确性
	count, err := page.Eval(`() => {
		if (window.__INITIAL_STATE__ && 
			window.__INITIAL_STATE__.search && 
			window.__INITIAL_STATE__.search.feeds && 
			window.__INITIAL_STATE__.search.feeds._value) {
			return window.__INITIAL_STATE__.search.feeds._value.length;
		}
		
		console.log('最终返回数量:', maxCount, '(state:', stateCount, ')');
		return maxCount;
	}`)

	if err != nil {
		fmt.Printf("获取帖子数量时出错: %v\n", err)
		return 0
	}

	if countInt := count.Value.Int(); countInt != 0 {
		return int(countInt)
	}

	return 0
}

func makeSearchURL(keyword string) string {

	values := url.Values{}
	values.Set("keyword", keyword)
	values.Set("source", "web_explore_feed")

	return fmt.Sprintf("https://www.xiaohongshu.com/search_result?%s", values.Encode())
}
