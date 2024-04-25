// Package logger provides a custom-configured zap.Logger.
package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevels is a struct that holds the log level and whether to add stacktrace or not.
type LogLevels struct {
	LogLevel      zapcore.Level
	AddStacktrace bool
}

// Logger is a custom-configured zap.Logger.
type Logger struct {
	*zap.Logger
	logLevel zap.AtomicLevel
}

// NewNop creates and returns a no-op zap.Logger for test
func NewNop() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
	}
}

// NewLogger creates and returns a custom-configured zap.Logger.
func NewLogger(lv *LogLevels) *Logger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	options := []zap.Option{
		zap.AddCaller(),
	}

	logLevel, stacktrace := getLogLevelFromEnvOrDefault(lv)
	if stacktrace {
		options = append(options, zap.AddStacktrace(logLevel))
	}

	atomicLevel := zap.NewAtomicLevelAt(logLevel)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	)

	logger := zap.New(core, options...)

	return &Logger{
		Logger:   logger,
		logLevel: atomicLevel,
	}
}

func getLogLevelFromEnvOrDefault(lv *LogLevels) (zapcore.Level, bool) {
	if lv != nil {
		return lv.LogLevel, lv.AddStacktrace
	}

	envLevel := zapcore.InfoLevel

	envLevelStr := os.Getenv("LOG_LEVEL")
	if envLevelStr != "" {
		err := envLevel.UnmarshalText([]byte(envLevelStr))
		if err != nil {
			log.Printf("Invalid LOG_LEVEL '%s', using default InfoLevel. Error: %v", envLevelStr, err)
		}
	}

	envStacktrace := false

	envStacktraceBool := os.Getenv("STACKTRACE")
	if envStacktraceBool == "true" {
		envStacktrace = true
	}

	return envLevel, envStacktrace
}

// SetLogLevel sets the log level of the logger.
func (l *Logger) SetLogLevel(level string) {
	l.Info("current log level", zap.String("level", l.logLevel.String()))

	err := l.logLevel.UnmarshalText([]byte(level))
	if err != nil {
		l.Error("invalid log level provided and continue with existing log level", zap.Error(err))
	}

	l.Info("new log level", zap.String("level", l.logLevel.String()))
}
