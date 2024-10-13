package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/henrywhitaker3/connect-template/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	statusOk         = 0
	statusBadRequest = 400
	statusUnauth     = 401
	statusBad        = 500
	statusTimeout    = 503
)

func Metrics() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				return next(ctx, req)
			}

			start := time.Now()
			info := strings.Split(strings.TrimLeft(req.Spec().Procedure, "/"), "/")
			service := info[0]
			method := info[1]
			protocol := req.Peer().Protocol
			resp, err := next(ctx, req)

			dur := time.Since(start)

			status := statusOk
			if err != nil {
				status = statusBad
				connerr := &connect.Error{}
				if errors.As(err, &connerr) {
					switch connerr.Code() {
					case connect.CodeCanceled:
						fallthrough
					case connect.CodeAborted:
						fallthrough
					case connect.CodeAlreadyExists:
						status = statusBadRequest
					case connect.CodeDeadlineExceeded:
						status = statusTimeout
					case connect.CodeUnauthenticated:
						status = statusUnauth
					}
				}
			}

			labels := prometheus.Labels{
				"service":  service,
				"method":   method,
				"protocol": protocol,
				"status":   fmt.Sprintf("%d", status),
			}

			metrics.GrpcRequestSeconds.With(labels).Observe(float64(dur) / float64(time.Second))

			return resp, err
		})
	}
}
