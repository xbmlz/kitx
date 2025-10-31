package log

import (
	"sync"

	"go.uber.org/zap"
)

// 全局 logger 实例
var (
	globalLogger *zap.Logger
	once         sync.Once
)

func init() {
	once.Do(func() {
		globalLogger = NewZapLogger()
	})
}

func Setup(opts ...Option) {
	globalLogger = NewZapLogger(opts...)
}

func GetLogger() *zap.Logger {
	return globalLogger
}

func Debug(format string, args ...any) {
	globalLogger.Sugar().Debugf(format, args...)
}

func Info(format string, args ...any) {
	globalLogger.Sugar().Infof(format, args...)
}

func Warn(format string, args ...any) {
	globalLogger.Sugar().Warnf(format, args...)
}

func Error(format string, args ...any) {
	globalLogger.Sugar().Errorf(format, args...)
}

func Fatal(format string, args ...any) {
	globalLogger.Sugar().Fatalf(format, args...)
}
