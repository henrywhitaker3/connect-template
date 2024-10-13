package middleware

import (
	"context"

	"connectrpc.com/connect"
	"github.com/henrywhitaker3/connect-template/internal/logger"
	"go.uber.org/zap"
)

func Zap(level zap.AtomicLevel) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx = logger.Wrap(ctx, level)
			return next(ctx, req)
		})
	}
}
