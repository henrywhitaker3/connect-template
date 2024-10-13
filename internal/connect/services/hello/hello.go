package hello

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	hellov1 "github.com/henrywhitaker3/connect-template/gen/hello/v1"
	"github.com/henrywhitaker3/connect-template/internal/app"
	"github.com/henrywhitaker3/connect-template/internal/logger"
)

type HelloServer struct {
	app *app.App
}

func New(app *app.App) *HelloServer {
	return &HelloServer{app: app}
}

func (h *HelloServer) Hello(
	ctx context.Context,
	req *connect.Request[hellov1.HelloRequest],
) (*connect.Response[hellov1.HelloResponse], error) {
	logger.Logger(ctx).Info("got hello")
	return connect.NewResponse(&hellov1.HelloResponse{
		Message: fmt.Sprintf("Hello %s!", req.Msg),
	}), nil
}
