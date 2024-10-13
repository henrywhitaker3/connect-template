package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/henrywhitaker3/connect-template/gen/hello/v1/hellov1connect"
	"github.com/henrywhitaker3/connect-template/internal/app"
	"github.com/stretchr/testify/require"
)

func Authenticate[T any](t *testing.T, app *app.App, req *connect.Request[T]) *connect.Request[T] {
	token, err := app.Jwt.New(time.Minute)
	require.Nil(t, err)
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req
}

func HelloClient(t *testing.T, app *app.App) (hellov1connect.HelloServiceClient, context.CancelFunc) {
	server := httptest.NewServer(app.Http)

	client := hellov1connect.NewHelloServiceClient(
		http.DefaultClient,
		server.URL,
	)

	return client, func() { server.Close() }
}

func ErrorIs(t *testing.T, err error, target error) {
	conn := &connect.Error{}
	require.ErrorAs(t, err, &conn)
	require.ErrorAs(t, conn.Unwrap(), &target)
}
