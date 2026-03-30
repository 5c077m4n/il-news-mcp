package main

import (
	"github.com/5c077m4n/il-news-mcp/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
