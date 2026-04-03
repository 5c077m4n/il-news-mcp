package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func New() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			start := time.Now()
			sessionID := req.GetSession().ID()

			slog.Info("REQUEST", "session", sessionID, "method", method)
			result, err := next(ctx, method, req)
			duration := time.Since(start)

			if err != nil {
				slog.Info(
					"RESPONSE",
					"session", sessionID,
					"method", method,
					"duration", duration,
					"error", err,
				)
			} else {
				slog.Info("RESPONSE", "session", sessionID, "method", method, "duration", duration)
			}

			return result, err
		}
	}
}
