package logger

import (
	"charonoms/internal/infrastructure/config"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

// Init 初始化日志系统
func Init(cfg config.LoggerConfig) error {
	var zapConfig zap.Config

	// 根据日志级别选择配置
	switch cfg.Level {
	case "debug":
		zapConfig = zap.NewDevelopmentConfig()
	case "info", "warn", "error":
		zapConfig = zap.NewProductionConfig()
	default:
		zapConfig = zap.NewProductionConfig()
	}

	// 设置日志级别
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// 设置输出
	if cfg.Output == "file" && cfg.FilePath != "" {
		// 确保日志目录存在
		if err := os.MkdirAll("logs", 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
		zapConfig.OutputPaths = []string{cfg.FilePath}
		zapConfig.ErrorOutputPaths = []string{cfg.FilePath}
	} else {
		zapConfig.OutputPaths = []string{"stdout"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}
	}

	// 自定义编码器配置
	zapConfig.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 构建日志实例
	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	Logger = logger
	SugarLogger = logger.Sugar()

	return nil
}

// Sync 同步日志
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
	if SugarLogger != nil {
		_ = SugarLogger.Sync()
	}
}

// Debug 输出 Debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 输出 Info 级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn 输出 Warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 输出 Error 级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 输出 Fatal 级别日志并退出
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Debugf 格式化输出 Debug 级别日志
func Debugf(template string, args ...interface{}) {
	SugarLogger.Debugf(template, args...)
}

// Infof 格式化输出 Info 级别日志
func Infof(template string, args ...interface{}) {
	SugarLogger.Infof(template, args...)
}

// Warnf 格式化输出 Warn 级别日志
func Warnf(template string, args ...interface{}) {
	SugarLogger.Warnf(template, args...)
}

// Errorf 格式化输出 Error 级别日志
func Errorf(template string, args ...interface{}) {
	SugarLogger.Errorf(template, args...)
}

// Fatalf 格式化输出 Fatal 级别日志并退出
func Fatalf(template string, args ...interface{}) {
	SugarLogger.Fatalf(template, args...)
}
