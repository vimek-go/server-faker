package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugMode  = "debug"
	normalMode = "normal"
)

type Logger interface {
	Sync() error
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
}

func NewExternalLogger(env string) (Logger, error) {
	logger, err := newConfig(env).Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func newConfig(env string) *zap.Config {
	config := zap.NewProductionConfig()
	if env == debugMode {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
	}
	config.Encoding = "console"
	config.EncoderConfig = newEncoderConfig(env)

	return &config
}

func newEncoderConfig(env string) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	if env == debugMode {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "severity"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeLevel = encodeLevel()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return encoderConfig
}

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		case zapcore.InvalidLevel:
			enc.AppendString("INFO")
		}
	}
}
