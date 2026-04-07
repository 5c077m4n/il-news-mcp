package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestServerResponse(t *testing.T) {
	server := createServer()
	handler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return server },
		nil,
	)
	testServer := httptest.NewServer(corsMiddleware(handler))
	defer testServer.Close()

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.1.0"}, nil)
	session, err := client.Connect(
		ctx,
		&mcp.StreamableClientTransport{Endpoint: testServer.URL, MaxRetries: -1},
		nil,
	)
	if assert.NoError(t, err) {
		defer assert.NoError(t, session.Close())
		assert.NotEmpty(t, session.ID())

		toolsResult, err := session.ListTools(ctx, nil)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, toolsResult.Tools)
			assert.Equal(t, "news", toolsResult.Tools[0].Name)

			callResult, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name:      toolsResult.Tools[0].Name,
				Arguments: map[string]any{},
			})
			assert.NoError(t, err)
			assert.NotNil(t, callResult)
			assert.NotEmpty(t, callResult.Content)
			assert.Greater(t, len(callResult.Content), 2)
		}
	}
}
