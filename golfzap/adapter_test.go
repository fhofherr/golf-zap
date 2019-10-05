// nolint: funlen
package golfzap_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/fhofherr/golf-zap/golfzap"
	"github.com/fhofherr/golf/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLog(t *testing.T) {
	tests := []struct {
		name     string
		kvs      []interface{}
		expected string
	}{
		{
			name:     "no level - no message",
			kvs:      []interface{}{"key", "value"},
			expected: `{"level":"info","msg":"no message","key":"value"}`,
		},
		{
			name:     "level debug - no message",
			kvs:      []interface{}{"level", "debug", "key", "value"},
			expected: `{"level":"debug","msg":"no message","key":"value"}`,
		},
		{
			name:     "lvl debug - no message",
			kvs:      []interface{}{"lvl", "debug", "key", "value"},
			expected: `{"level":"debug","msg":"no message","key":"value"}`,
		},
		{
			name:     "level info - no message",
			kvs:      []interface{}{"level", "info", "key", "value"},
			expected: `{"level":"info","msg":"no message","key":"value"}`,
		},
		{
			name:     "lvl info - no message",
			kvs:      []interface{}{"lvl", "info", "key", "value"},
			expected: `{"level":"info","msg":"no message","key":"value"}`,
		},
		{
			name:     "level warn - no message",
			kvs:      []interface{}{"level", "warn", "key", "value"},
			expected: `{"level":"warn","msg":"no message","key":"value"}`,
		},
		{
			name:     "lvl warn - no message",
			kvs:      []interface{}{"lvl", "warn", "key", "value"},
			expected: `{"level":"warn","msg":"no message","key":"value"}`,
		},
		{
			name:     "level error - no message",
			kvs:      []interface{}{"level", "error", "key", "value"},
			expected: `{"level":"error","msg":"no message","key":"value"}`,
		},
		{
			name:     "lvl error - no message",
			kvs:      []interface{}{"lvl", "error", "key", "value"},
			expected: `{"level":"error","msg":"no message","key":"value"}`,
		},
		{
			name:     "no level - some message",
			kvs:      []interface{}{"message", "some message", "key", "value"},
			expected: `{"level":"info","msg":"some message","key":"value"}`,
		},
		{
			name:     "no level - some msg",
			kvs:      []interface{}{"msg", "some message", "key", "value"},
			expected: `{"level":"info","msg":"some message","key":"value"}`,
		},
		{
			name:     "some level - some message",
			kvs:      []interface{}{"level", "debug", "message", "some message", "key", "value"},
			expected: `{"level":"debug","msg":"some message","key":"value"}`,
		},
		{
			name:     "some lvl - some msg",
			kvs:      []interface{}{"lvl", "debug", "msg", "some message", "key", "value"},
			expected: `{"level":"debug","msg":"some message","key":"value"}`,
		},
		{
			name:     "level and lvl",
			kvs:      []interface{}{"lvl", "debug", "level", "warn"},
			expected: `{"level":"warn","msg":"no message"}`,
		},
		{
			name:     "later level wins",
			kvs:      []interface{}{"level", "debug", "level", "warn"},
			expected: `{"level":"warn","msg":"no message"}`,
		},
		{
			name:     "message and msg",
			kvs:      []interface{}{"message", "the message", "msg", "some msg"},
			expected: `{"level":"info","msg":"the message"}`,
		},
		{
			name:     "invalid level type",
			kvs:      []interface{}{"level", 1},
			expected: `{"level":"info","msg":"no message"}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			wsa := &writerSyncerAdapter{writer: buf}
			logger := newZapAdapter(wsa)
			err := logger.Log(tt.kvs...)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, buf.String())
		})
	}
}

func TestWith(t *testing.T) {
	tests := []struct {
		name     string
		withKVs  []interface{}
		kvs      []interface{}
		expected string
	}{
		{
			name:     "some kvs and message",
			withKVs:  []interface{}{"key1", "value1"},
			kvs:      []interface{}{"message", "some message"},
			expected: `{"level": "info", "key1": "value1", "msg": "some message"}`,
		},
		{
			name:     "later level does not override prior level",
			withKVs:  []interface{}{"level", "debug"},
			kvs:      []interface{}{"level", "error"},
			expected: `{"level": "debug", "msg": "no message"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			wsa := &writerSyncerAdapter{writer: buf}
			logger := newZapAdapter(wsa)
			logger = log.With(logger, tt.withKVs...)
			err := logger.Log(tt.kvs...)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, buf.String())
		})
	}
}

func newZapAdapter(ws zapcore.WriteSyncer) log.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), ws, zapcore.DebugLevel)
	logger := zap.New(core)
	return golfzap.New(logger)
}

type writerSyncerAdapter struct {
	writer io.Writer
}

func (wsa writerSyncerAdapter) Write(p []byte) (n int, err error) {
	return wsa.writer.Write(p)
}

func (wsa writerSyncerAdapter) Sync() error {
	// no-op
	return nil
}
