package golfzap

import (
	"github.com/fhofherr/golf/log"
	"go.uber.org/zap"
)

const (
	defaultLogLevel   = "info"
	defaultLogMessage = "no message"
)

type adapter struct {
	Logger *zap.Logger
}

// New creates a new golf adapter wrapping the passed Logger.
func New(logger *zap.Logger) log.Logger {
	return &adapter{
		Logger: logger,
	}
}

// Log logs the passed key-value pairs using the wrapped zap.SugaredLogger.
//
// If kvs contains the keys "lvl" or "level" the corresponding value will be
// be used to determine the log level. The following values are possible:
//
// - "debug"
// - "info"
// - "warn"
// - "error"
//
// Log calls the corresponding method on the wrapped zap.SugaredLogger, i.e. if
// "level" or "lvl" has the value "warn" zap.SugaredLogger.Warn will be used
// to log the message. Log removes the "level" or "lvl" key along with its
// corresponding value before it calls the respective method on
// zap.SugaredLogger.
//
// If the value of "lvl" or "level" does not match one of the above, or if kvs
// does not contain "lvl" or "level" at all, the info level is assumed.
//
// If kvs contains both "lvl" and "level" Log gives preference to the "level"
// key.
//
// Likewise, if kvs contains the key "msg" or "message" Log will treat its value
// as a log message. As with the logging level Log removes the "msg" or "message"
// key and its corresponding value from kvs, before using the remainder as
// key-value pairs for zap.
//
// If both keys "msg" and "message" are present in kvs Log gives preference to
// "message".
//
// If neither "msg" nor "message" are present, Log passes the value "no message"
// as message to the respective zap method.
func (a *adapter) Log(kvs ...interface{}) error {
	logger := a.Logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	level, msg, entry := prepareEntry(kvs)
	switch level {
	case "debug":
		logger.Debugw(msg, entry...)
	case "info":
		logger.Infow(msg, entry...)
	case "warn":
		logger.Warnw(msg, entry...)
	case "error":
		logger.Errorw(msg, entry...)
	}
	return nil
}

func prepareEntry(kvs []interface{}) (string, string, []interface{}) {
	entry := make([]interface{}, 0, len(kvs))
	level := defaultLogLevel
	msg := defaultLogMessage
	for i := 0; i < len(kvs); {
		if isSpecialKey(kvs, i, "level", "lvl") {
			if kvs[i] == "level" || level == defaultLogLevel {
				level = toString(kvs[i+1], defaultLogLevel)
			}
			i += 2
			continue
		}
		if isSpecialKey(kvs, i, "message", "msg") {
			if kvs[i] == "message" || msg == defaultLogMessage {
				msg = toString(kvs[i+1], defaultLogMessage)
			}
			i += 2
			continue
		}
		entry = append(entry, kvs[i])
		i++
	}
	return level, msg, entry
}

func isSpecialKey(kvs []interface{}, i int, specialKeys ...string) bool {
	if i+1 >= len(kvs) {
		return false
	}
	for _, k := range specialKeys {
		if k == kvs[i] {
			return true
		}
	}
	return false
}

func toString(x interface{}, d string) string {
	if s, ok := x.(string); ok {
		return s
	}
	return d
}

// With adds the values in kvs to the zap adapter.
//
// With delegates to the wrapped zap logger's With method. It therefore depends
// on how zap handles key value pairs passed to With. Especially this means
// that passing "level" to with defines the log level for this logger. It is not
// possible to override a "level" key passed to With by passing a different level
// to Log.
func (a *adapter) With(kvs ...interface{}) log.Logger {
	logger := a.Logger.Sugar()
	logger = logger.With(kvs...)
	return &adapter{
		Logger: logger.Desugar(),
	}
}
