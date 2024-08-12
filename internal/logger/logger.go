package logger

import (
	"context"

	"github.com/henrywhitaker3/ctxgen"
	"go.uber.org/zap"
)

var (
	zl *zap.SugaredLogger
)

func Wrap(ctx context.Context, level zap.AtomicLevel) context.Context {
	return ctxgen.WithValue(ctx, "logger", newLogger(level))
}

func Logger(ctx context.Context) *zap.SugaredLogger {
	log, ok := ctxgen.ValueOk[*zap.SugaredLogger](ctx, "logger")
	if !ok {
		log = newLogger(zap.NewAtomicLevelAt(zap.ErrorLevel))
	}
	return log
}

func newLogger(level zap.AtomicLevel) *zap.SugaredLogger {
	if zl == nil {
		logConfig := zap.NewProductionConfig()
		logConfig.OutputPaths = []string{"stdout"}
		logConfig.Level = level
		logger, _ := logConfig.Build()
		zl = logger.Sugar()
	}
	return zl
}
