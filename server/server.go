package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/5c077m4n/il-news-mcp/server/feed"
	"github.com/5c077m4n/il-news-mcp/server/logger"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	version = "0.1.0"
	port    = "8888"
	url     = "127.0.0.1:" + port
)

type getNewsParams struct {
	DateRange [2]time.Time `json:"dateRange" jsonschema:"start/end dates for the news"`
}

func getNews(
	ctx context.Context,
	_req *mcp.CallToolRequest,
	_params *getNewsParams,
) (*mcp.CallToolResult, any, error) {
	feedAgg := sync.Map{}
	var wg sync.WaitGroup

	wg.Go(func() {
		ynetFeed, err := feed.GetYnet(ctx)
		if err != nil {
			slog.Error("could not fetch Ynet feed", "error", err.Error())
			return
		}

		feedAgg.Store("ynet", ynetFeed)
	})
	wg.Go(func() {
		abuFeed, err := feed.GetAbuAliExpress(ctx)
		if err != nil {
			slog.Error("could not fetch Abu Ali Express feed", "error", err.Error())
			return
		}

		feedAgg.Store("abu_ali_express", abuFeed)
	})
	wg.Wait()

	tmpFeedAgg := make(map[string]any)
	feedAgg.Range(func(key, value any) bool {
		tmpFeedAgg[key.(string)] = value
		return true
	})
	data, err := json.Marshal(tmpFeedAgg)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func Run() error {
	server := mcp.NewServer(&mcp.Implementation{Name: "il-news-mcp", Version: version}, nil)
	server.AddReceivingMiddleware(logger.New())
	mcp.AddTool(server, &mcp.Tool{
		Name:        "news",
		Description: "Get the most relevant news",
	}, getNews)

	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		url := request.URL.Path
		switch url {
		default:
			return server
		}
	}, nil)

	slog.Info("MCP server listening", "URL", url)
	return http.ListenAndServe(url, handler)
}
