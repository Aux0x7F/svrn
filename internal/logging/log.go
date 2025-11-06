package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		fallback := zap.NewExample().Sugar()
		fallback.Error("failed to initialize zap; using fallback", "error", err)
		return fallback
	}
	return logger.Sugar()
}

func ReplaceGlobals(l *zap.SugaredLogger) {
	zap.ReplaceGlobals(l.Desugar())
}

func Sync(l *zap.SugaredLogger) {
	_ = l.Sync()
}
