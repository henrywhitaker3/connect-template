package connect

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	connectrpc "connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/henrywhitaker3/connect-template/gen/hello/v1/hellov1connect"
	"github.com/henrywhitaker3/connect-template/internal/app"
	"github.com/henrywhitaker3/connect-template/internal/connect/middleware"
	"github.com/henrywhitaker3/connect-template/internal/connect/services/hello"
	"github.com/henrywhitaker3/connect-template/internal/logger"
	"github.com/henrywhitaker3/connect-template/internal/tracing"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Registerer func() (string, http.Handler)

type Server struct {
	app    *app.App
	server *http.Server
	mux    *http.ServeMux
}

func New(app *app.App) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		app: app,
		mux: mux,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", app.Config.Http.Port),
			Handler: h2c.NewHandler(mux, &http2.Server{}),
		},
	}

	opts := []otelconnect.Option{}
	if app.Config.Telemetry.Tracing.Enabled {
		opts = append(opts, otelconnect.WithTracerProvider(tracing.TracerProvider))
	}
	ot, err := otelconnect.NewInterceptor(opts...)
	if err != nil {
		panic(err)
	}
	base := connectrpc.WithInterceptors(
		middleware.Errors(),
		middleware.Zap(app.Config.LogLevel.Level()),
		ot,
		middleware.Metrics(),
	)
	auth := connectrpc.WithInterceptors(middleware.NewAuth(app.Jwt))

	srv.Register(hellov1connect.NewHelloServiceHandler(
		hello.New(app),
		base,
		auth,
	))

	return srv
}

func (s *Server) Register(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

func (s *Server) Start(ctx context.Context) error {
	logger.Logger(ctx).Infow("starting connect server", "port", s.app.Config.Http.Port)
	if err := s.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Logger(ctx).Info("stopping connect server")
	return s.server.Shutdown(ctx)
}
