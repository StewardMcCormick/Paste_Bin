package util

import (
	"context"
	"github.com/StewardMcCormick/Paste_Bin/internal/util/http_util"
	"go.uber.org/zap"
)

func GetLoggerFromCtx(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(http_util.LoggerKey).(*zap.Logger); ok {
		return logger
	}
	return zap.L()
}
