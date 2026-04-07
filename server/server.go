package server

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/5c077m4n/il-news-mcp/server/middleware/cors"
	"github.com/5c077m4n/il-news-mcp/server/middleware/logger"
	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

var corsMiddleware = cors.New()

func createServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "il-news-mcp", Version: version},
		&mcp.ServerOptions{
			Logger:       slog.Default().WithGroup("mcp_server"),
			GetSessionID: func() string { return uuid.New().String() },
		},
	)
	server.AddReceivingMiddleware(logger.New())
	mcp.AddTool(server, &mcp.Tool{Name: "news", Description: "Get the most relevant news"}, getNews)

	return server
}

func parseCLIArgs() (string, uint) {
	host := flag.String("host", "localhost", "the host address to run this server on")
	port := flag.Uint("port", 8888, "the port to run this server on")

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

	return *host, *port
}

func Run() error {
	host, port := parseCLIArgs()

	server := createServer()
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		path := request.URL.Path
		slog.Info("Handling request", "pathname", path)

		return server
	}, nil)

	serverURL := fmt.Sprintf("%s:%d", host, port)
	slog.Info("MCP server listening", "URL", serverURL)

	return http.ListenAndServe(serverURL, corsMiddleware(handler))
}
