package log

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(opts ...Option) *zap.Logger {
	cfg := Config{
		Level:      "info",
		Format:     "console",
		Filename:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	cores := []zapcore.Core{
		createConsoleZapCore(cfg.Level, cfg.Format),
	}

	if cfg.Filename != "" {
		cores = append(cores, createFileZapCore(cfg))
	}

	return zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
}

func ParseLevel(level string) zapcore.Level {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return zapcore.InfoLevel
	}
	return l
}

func createConsoleZapCore(level, encoder string) zapcore.Core {
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	consoleEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	consoleEncoderConfig.EncodeName = zapcore.FullNameEncoder
	consoleEncoderConfig.ConsoleSeparator = "\t"
	return zapcore.NewCore(
		createEncoderConfig(encoder, consoleEncoderConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zap.NewAtomicLevelAt(ParseLevel(level)),
	)
}

func createFileZapCore(cfg Config) zapcore.Core {
	if err := os.MkdirAll(filepath.Dir(cfg.Filename), os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	fileEncoderConfig.EncodeName = zapcore.FullNameEncoder
	fileEncoderConfig.ConsoleSeparator = "\t"
	hook := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.NewCore(
		createEncoderConfig(cfg.Format, fileEncoderConfig),
		zapcore.AddSync(hook),
		zap.NewAtomicLevelAt(ParseLevel(cfg.Level)),
	)
}

func createEncoderConfig(encoder string, consoleEncoderConfig zapcore.EncoderConfig) zapcore.Encoder {
	if encoder == "json" {
		return zapcore.NewJSONEncoder(consoleEncoderConfig)
	}
	return zapcore.NewConsoleEncoder(consoleEncoderConfig)
}
