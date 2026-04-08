package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/5c077m4n/il-news-mcp/server/feed"
	"github.com/5c077m4n/pikud-haoref-api-go/history"
	"github.com/goccy/go-json"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type getNewsParams struct {
	right *bool
	left  *bool
}

func getNews(
	ctx context.Context,
	_req *mcp.CallToolRequest,
	params getNewsParams,
) (*mcp.CallToolResult, any, error) {
	content := NewSafeList[mcp.Content]()
	wg := sync.WaitGroup{}

	for source, getterFn := range feed.NewsSourceToGetter {
		if orientation, found := feed.NewsSourceToOrientation[source]; params.right != nil &&
			!*params.right &&
			found &&
			orientation > 0 {
			continue
		}
		if orientation, found := feed.NewsSourceToOrientation[source]; params.left != nil &&
			!*params.left &&
			found &&
			orientation < 0 {
			continue
		}

		wg.Go(func() {
			feedContent, err := getterFn(ctx)
			if err != nil {
				slog.WarnContext(ctx, "could not fetch feed", "source", source, "error", err)
				return
			}

			feedContentBytes, err := json.MarshalContext(ctx, feedContent)
			if err != nil {
				slog.WarnContext(
					ctx,
					"could not marshal content",
					"error",
					err,
					"contents",
					feedContent,
				)
				return
			}

			orientation := feed.NewsSourceToOrientation[source]
			content.Append(&mcp.TextContent{
				Text: string(feedContentBytes),
				Meta: mcp.Meta{
					"fetchedAt":   time.Now(),
					"source":      source,
					"orientation": orientation,
				},
			},
			)
		})
	}
	wg.Wait()

	return &mcp.CallToolResult{Content: content.ToList()}, nil, nil
}

func getMissileAlerts(
	ctx context.Context,
	_req *mcp.CallToolRequest,
	_params struct{},
) (*mcp.CallToolResult, any, error) {
	alerts, err := history.FetchAlerts(ctx)
	if err != nil {
		return nil, nil, err
	}

	alertsBytes, err := json.MarshalContext(ctx, alerts)
	if err != nil {
		return nil, nil, err
	}

	content := []mcp.Content{
		&mcp.TextContent{
			Text: string(alertsBytes),
			Meta: mcp.Meta{"fetchedAt": time.Now()},
		},
	}
	return &mcp.CallToolResult{Content: content}, nil, nil
}
