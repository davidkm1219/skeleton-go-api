package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/twk/skeleton-go-api/internal/logger"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		lv  *logger.LogLevels
		env map[string]string
	}

	type want struct {
		level zapcore.Level
	}

	tests := map[string]struct {
		args args
		want want
	}{
		"Default": {
			args: args{
				lv: &logger.LogLevels{},
			},
			want: want{
				level: zapcore.InfoLevel,
			},
		},
		"Debug": {
			args: args{
				lv: &logger.LogLevels{
					LogLevel: zapcore.DebugLevel,
				},
			},
			want: want{
				level: zapcore.DebugLevel,
			},
		},
		"Log level from env": {
			args: args{
				lv: nil,
				env: map[string]string{
					"LOG_LEVEL": "debug",
				},
			},
			want: want{
				level: zapcore.DebugLevel,
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			for k, v := range tt.args.env {
				t.Setenv(k, v)
			}

			l := logger.NewLogger(tt.args.lv)
			assert.Equal(t, tt.want.level, l.Level())
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	lv := &logger.LogLevels{
		LogLevel:      zapcore.DebugLevel,
		AddStacktrace: true,
	}
	l := logger.NewLogger(lv)
	observedZapCore, logs := observer.New(zap.DebugLevel)
	l.Logger = l.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core { return observedZapCore }))

	l.Debug("test message")

	if logs.Len() != 1 {
		t.Errorf("Expected 1 log message, got %d", logs.Len())
	}

	if logs.All()[0].Message != "test message" {
		t.Errorf("Expected 'test message', got '%s'", logs.All()[0].Message)
	}
}
