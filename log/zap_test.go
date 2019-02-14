package log

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"testing"

	"go.uber.org/zap"
)

// type

func TestZap(t *testing.T) {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "x.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zap.DebugLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	for i := 0; i < 10; i++ {
		logger.Info("xxx", zap.Int("index", i))
	}
	logger.Error("mmmmmmmm", zap.String("err", "dasdasdsdsa"))
}
