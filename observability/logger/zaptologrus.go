package logger

import (
	"runtime"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogrusZapHook struct {
	Logger *zap.Logger
}

func NewLogrusZapHook(logger *zap.Logger) (*LogrusZapHook, error) {
	return &LogrusZapHook{
		Logger: logger,
	}, nil
}

func (hook *LogrusZapHook) Fire(entry *logrus.Entry) error {
	fields := make([]zap.Field, 0, 10)

	for key, value := range entry.Data {
		if key == logrus.ErrorKey {
			fields = append(fields, zap.Error(value.(error)))
		} else {
			fields = append(fields, zap.Any(key, value))
		}
	}

	switch entry.Level {
	case logrus.PanicLevel:
		hook.Write(zapcore.PanicLevel, entry.Message, fields, entry.Caller)
	case logrus.FatalLevel:
		hook.Write(zapcore.FatalLevel, entry.Message, fields, entry.Caller)
	case logrus.ErrorLevel:
		hook.Write(zapcore.ErrorLevel, entry.Message, fields, entry.Caller)
	case logrus.WarnLevel:
		hook.Write(zapcore.WarnLevel, entry.Message, fields, entry.Caller)
	case logrus.InfoLevel:
		hook.Write(zapcore.InfoLevel, entry.Message, fields, entry.Caller)
	case logrus.DebugLevel, logrus.TraceLevel:
		hook.Write(zapcore.DebugLevel, entry.Message, fields, entry.Caller)
	}

	return nil
}

func (hook *LogrusZapHook) Write(lvl zapcore.Level, msg string, fields []zap.Field, caller *runtime.Frame) {
	if ce := hook.Logger.Check(lvl, msg); ce != nil {
		if caller != nil {
			ce.Caller = zapcore.NewEntryCaller(caller.PC, caller.File, caller.Line, caller.PC != 0)
		}
		ce.Write(fields...)
	}
}

func (hook *LogrusZapHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
