package logutil

import (
	"context"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFilePath      = "./log/server.log"
	errorLogFilePath = "./log/error.log"
	logMaxSize       = 128
	logMaxCount      = 30
	logMaxAge        = 7
)

var logger *zap.SugaredLogger

func InitLog() {
	defaultEncoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	logWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    logMaxSize,
		MaxAge:     logMaxAge,
		MaxBackups: logMaxCount,
		Compress:   false,
	})
	errorLogWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFilePath,
		MaxSize:    logMaxSize,
		MaxAge:     logMaxAge,
		MaxBackups: logMaxCount,
		Compress:   false,
	})
	consolecore := zapcore.NewConsoleEncoder(defaultEncoder)

	fileEncoder := defaultEncoder
	fileEncoder.EncodeLevel = zapcore.CapitalLevelEncoder
	fileCore := zapcore.NewConsoleEncoder(fileEncoder)

	core := zapcore.NewTee(
		zapcore.NewCore(fileCore, logWriteSyncer, zap.InfoLevel),
		zapcore.NewCore(fileCore, errorLogWriteSyncer, zap.ErrorLevel),
		zapcore.NewCore(consolecore, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func LoggerWithContext(ctx context.Context) *zap.SugaredLogger {
	return logger.With()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
