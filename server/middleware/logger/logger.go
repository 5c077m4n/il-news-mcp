package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var logger = slog.Default().WithGroup("mcp_server_middleware")

func New() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			start := time.Now()
			sessionID := req.GetSession().ID()

			logger.Info("REQUEST", "session_id", sessionID, "method", method)
			result, err := next(ctx, method, req)
			duration := time.Since(start)

			if err != nil {
				logger.Info(
					"RESPONSE",
					"session_id", sessionID,
					"method", method,
					"duration", duration,
					"error", err,
				)
			} else {
				logger.Info("RESPONSE", "session_id", sessionID, "method", method, "duration", duration)
			}

			return result, err
		}
	}
}
