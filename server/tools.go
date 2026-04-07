package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/5c077m4n/il-news-mcp/server/feed"
	"github.com/goccy/go-json"
	"github.com/mmcdole/gofeed"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type getNewsParams struct {
	right *bool
	left  *bool
}

func getNews(
	ctx context.Context,
	_req *mcp.CallToolRequest,
	_params struct{},
) (*mcp.CallToolResult, any, error) {
	feedAgg := sync.Map{}
	wg := sync.WaitGroup{}

	for source, getterFn := range feed.NewsSourceToGetter {
		wg.Go(func() {
			feedContent, err := getterFn(ctx)
			if err != nil {
				slog.WarnContext(ctx, "could not fetch feed", "source", source, "error", err)
				return
			}

			feedAgg.Store(source, feedContent.Items)
		})
	}
	wg.Wait()

	content := []mcp.Content{}
	for source, orientation := range feed.NewsSourceToOrientation {
		if feedContent, found := feedAgg.Load(source); found {
			feedContent := feedContent.([]*gofeed.Item)

			if feedContentBytes, err := json.MarshalContext(ctx, feedContent); err == nil {
				content = append(
					content,
					&mcp.TextContent{
						Text: string(feedContentBytes),
						Meta: mcp.Meta{
							"fetchedAt":   time.Now(),
							"source":      source,
							"orientation": orientation,
						},
					},
				)
			} else {
				slog.WarnContext(ctx, "could not fetch YNet's RSS feed", "error", err)
			}
		}
	}

	return &mcp.CallToolResult{Content: content}, nil, nil
}
