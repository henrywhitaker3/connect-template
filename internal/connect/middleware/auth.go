package middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/henrywhitaker3/connect-template/internal/connect/common"
	"github.com/henrywhitaker3/connect-template/internal/jwt"
)

func NewAuth(jwt *jwt.Jwt) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if !req.Spec().IsClient {
				token := req.Header().Get("Authorization")
				if !strings.Contains(token, "Bearer ") {
					return unauth()
				}
				token = strings.ReplaceAll(token, "Bearer ", "")
				if jwt.Verify(ctx, token) != nil {
					return unauth()
				}
			}
			return next(ctx, req)
		})
	}
}

func unauth() (connect.AnyResponse, error) {
	return nil, connect.NewError(
		connect.CodeUnauthenticated,
		common.ErrUnauthenticated,
	)
}
