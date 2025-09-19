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
		scrollInterval = 2 * time.Second
	}

	// 获取初始的帖子数量，用于监控是否有新内容加载
	initialCount := s.getCurrentFeedCount(page)

	fmt.Printf("开始滚动加载更多内容，初始帖子数量: %d\n", initialCount)

	for i := 0; i < scrollCount; i++ {
		fmt.Printf("第 %d/%d 次滚动...\n", i+1, scrollCount)

		// 执行滚动操作
		_, err := page.Eval(`() => {
			// 滚动到页面底部
			window.scrollTo({
				top: document.body.scrollHeight,
				behavior: 'smooth'
			});
			
			// 记录滚动信息
			console.log('已滚动到位置:', window.pageYOffset);
			return window.pageYOffset;
		}`)

		if err != nil {
			return fmt.Errorf("scroll operation failed: %w", err)
		}

		// 等待内容加载
		time.Sleep(scrollInterval)

		// 等待页面稳定
		page.MustWaitStable()

		// 检查是否有新内容加载
		currentCount := s.getCurrentFeedCount(page)
		fmt.Printf("滚动后帖子数量: %d (新增: %d)\n", currentCount, currentCount-initialCount)

		// 如果连续几次滚动都没有新内容，可能已经到底了
		if i > 0 && currentCount == initialCount {
			fmt.Println("检测到没有新内容加载，可能已到达页面底部")
			break
		}

		initialCount = currentCount
	}

	fmt.Printf("滚动完成，最终帖子数量: %d\n", initialCount)
	return nil
}

// getCurrentFeedCount 获取当前页面的帖子数量
func (s *SearchAction) getCurrentFeedCount(page *rod.Page) int {
	count, err := page.Eval(`() => {
		if (window.__INITIAL_STATE__ && 
			window.__INITIAL_STATE__.search && 
			window.__INITIAL_STATE__.search.feeds && 
			window.__INITIAL_STATE__.search.feeds._value) {
			return window.__INITIAL_STATE__.search.feeds._value.length;
		}
		return 0;
	}`)

	if err != nil {
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
