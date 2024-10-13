package middleware

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/getsentry/sentry-go"
	"github.com/henrywhitaker3/connect-template/internal/connect/common"
	"github.com/henrywhitaker3/connect-template/internal/logger"
)

func Errors() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)

			if err == nil {
				return resp, nil
			}

			switch true {
			case errors.Is(err, common.ErrUnauthenticated):
				return resp, connect.NewError(
					connect.CodeUnauthenticated,
					common.ErrUnauthenticated,
				)
			}

			logger.Logger(ctx).Errorw("unhandled error", "error", err)
			if hub := sentry.GetHubFromContext(ctx); hub != nil {
				hub.CaptureException(err)
			}

			return resp, connect.NewError(
				connect.CodeUnknown,
				errors.New("something went wrong"),
			)
		})
	}
}
