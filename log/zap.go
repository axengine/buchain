package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	infoOnly      struct{}
	infoWithDebug struct{}
	aboveWarn     struct{}
)

var Logger *zap.Logger

func (l infoOnly) Enabled(lv zapcore.Level) bool {
	return lv == zapcore.InfoLevel
}
func (l infoWithDebug) Enabled(lv zapcore.Level) bool {
	return lv == zapcore.InfoLevel || lv == zapcore.DebugLevel
}
func (l aboveWarn) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.WarnLevel
}

func makeInfoFilter(env string) zapcore.LevelEnabler {
	switch env {
	case "production":
		return infoOnly{}
	default:
		return infoWithDebug{}
	}
}

func makeErrorFilter() zapcore.LevelEnabler {
	return aboveWarn{}
}

func init() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/zap_debug.log",
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

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func Initialize(mode, env, output string) *zap.Logger {
	if mode == "file" {
		var encoder zapcore.Encoder
		if env == "production" {
			encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		} else {
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		}
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   output,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		})
		core := zapcore.NewCore(
			encoder,
			w,
			zap.DebugLevel,
		)

		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		var encoderCfg zapcore.EncoderConfig
		if env == "production" {
			encoderCfg = zap.NewProductionEncoderConfig()
		} else {
			encoderCfg = zap.NewDevelopmentEncoderConfig()
		}

		coreInfo := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(os.Stdout),
			makeInfoFilter(env),
		)
		coreError := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(os.Stderr),
			makeErrorFilter(),
		)

		Logger = zap.New(zapcore.NewTee(coreInfo, coreError))
	}

	return Logger
}
