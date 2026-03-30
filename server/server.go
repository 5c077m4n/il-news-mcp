package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
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
	_ctx context.Context,
	_req *mcp.CallToolRequest,
	_params *getNewsParams,
) (*mcp.CallToolResult, any, error) {
	feed, err := feed.GetYnet()
	if err != nil {
		return nil, nil, err
	}

	data, err := json.Marshal(feed)
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

	handler := mcp.NewStreamableHTTPHandler(func(_req *http.Request) *mcp.Server {
		return server
	}, nil)

	slog.Info("MCP server listening", "URL", url)
	return http.ListenAndServe(url, handler)
}
