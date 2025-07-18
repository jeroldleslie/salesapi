package log

import "go.uber.org/zap"

func InitLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// Only error console
	return logger
}
