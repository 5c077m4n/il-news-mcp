package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/5c077m4n/il-news-mcp/server/feed"
	"github.com/5c077m4n/il-news-mcp/server/middleware/cors"
	"github.com/5c077m4n/il-news-mcp/server/middleware/logger"
	"github.com/goccy/go-json"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

var corsMiddleware = cors.New()

func getNews(
	ctx context.Context,
	_req *mcp.CallToolRequest,
	_params struct{},
) (*mcp.CallToolResult, any, error) {
	feedAgg := sync.Map{}
	var wg sync.WaitGroup

	wg.Go(func() {
		ynetFeed, err := feed.GetYnet(ctx)
		if err != nil {
			slog.Error("could not fetch Ynet feed", "error", err)
			return
		}

		feedAgg.Store("ynet", ynetFeed.Channel.Items)
	})
	wg.Wait()

	content := []mcp.Content{}
	if ynetFeed, found := feedAgg.Load("ynet"); found {
		ynetFeed := ynetFeed.([]feed.YNetRSSItem)

		if ynetFeedBytes, err := json.MarshalContext(ctx, ynetFeed); err == nil {
			content = append(
				content,
				&mcp.TextContent{
					Text: string(ynetFeedBytes),
					Meta: mcp.Meta{"fetchedAt": time.Now(), "source": "ynet", "orientation": -5},
				},
			)
		} else {
			slog.Warn("could not fetch YNet's RSS feed", "error", err)
		}
	}

	return &mcp.CallToolResult{Content: content}, nil, nil
}

func Run() error {
	host := flag.String("host", "0.0.0.0", "the host address to run this server on")
	port := flag.Int("port", 8888, "the port to run this server on")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This program runs an MCP Israeli news server over SSE HTTP.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEndpoints:\n")
		fmt.Fprintf(os.Stderr, "\t/news - get the latest news\n")
		os.Exit(1)
	}
	flag.Parse()

	server := mcp.NewServer(
		&mcp.Implementation{Name: "il-news-mcp", Version: version},
		&mcp.ServerOptions{Logger: slog.Default().WithGroup("mcp_server")},
	)
	server.AddReceivingMiddleware(logger.New())
	mcp.AddTool(server, &mcp.Tool{Name: "news", Description: "Get the most relevant news"}, getNews)

	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
		url := request.URL.Path
		slog.Info("Handling request", "URL", url)

		switch url {
		default:
			return server
		}
	}, nil)

	serverURL := fmt.Sprintf("%s:%d", *host, *port)
	slog.Info("MCP server listening", "URL", serverURL)

	return http.ListenAndServe(serverURL, corsMiddleware(handler))
}
