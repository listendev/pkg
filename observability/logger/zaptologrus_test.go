package logger

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogrusZapHook(t *testing.T) {
	tests := []struct {
		name         string
		logrusLogger *logrus.Logger
		expected     []observer.LoggedEntry
		testFunc     func(t *testing.T, l *logrus.Logger, r *observer.ObservedLogs)
	}{
		{
			name:         "simple conversion",
			logrusLogger: logrus.New(),
			expected: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:      zapcore.InfoLevel,
						Time:       time.Time{},
						LoggerName: "",
						Message:    "here you go",
						Caller: zapcore.EntryCaller{
							Defined:  false,
							PC:       uintptr(0),
							File:     "",
							Line:     0,
							Function: "",
						},
						Stack: "",
					},
					Context: []zapcore.Field{
						{
							Key:       "component",
							Type:      zapcore.StringType,
							Integer:   0,
							String:    "dummyscheduler",
							Interface: nil,
						},
					},
				},
			},
			testFunc: func(t *testing.T, l *logrus.Logger, r *observer.ObservedLogs) {
				l.WithField("component", "dummyscheduler").Info("here you go")
			},
		},
		{
			name:         "multiline conversion",
			logrusLogger: logrus.New(),
			expected: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:      zapcore.InfoLevel,
						Time:       time.Time{},
						LoggerName: "",
						Message:    "here you go",
						Caller: zapcore.EntryCaller{
							Defined:  false,
							PC:       uintptr(0),
							File:     "",
							Line:     0,
							Function: "",
						},
						Stack: "",
					},
					Context: []zapcore.Field{
						{
							Key:       "component",
							Type:      zapcore.StringType,
							Integer:   0,
							String:    "dummyscheduler",
							Interface: nil,
						},
					},
				},
				{
					Entry: zapcore.Entry{
						Level:      zapcore.ErrorLevel,
						Time:       time.Time{},
						LoggerName: "",
						Message:    "what happened",
						Caller: zapcore.EntryCaller{
							Defined:  false,
							PC:       uintptr(0),
							File:     "",
							Line:     0,
							Function: "",
						},
						Stack: "",
					},
					Context: []zapcore.Field{
						{
							Key:       "component",
							Type:      zapcore.StringType,
							Integer:   0,
							String:    "dummyscheduler",
							Interface: nil,
						},
					},
				},
			},
			testFunc: func(t *testing.T, l *logrus.Logger, r *observer.ObservedLogs) {
				l.WithField("component", "dummyscheduler").Info("here you go")
				l.WithField("component", "dummyscheduler").Error("what happened")
			},
		},
		{
			name:         "multiline conversion",
			logrusLogger: logrus.New(),
			expected: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:      zapcore.InfoLevel,
						Time:       time.Time{},
						LoggerName: "",
						Message:    "here you go",
						Caller: zapcore.EntryCaller{
							Defined:  false,
							PC:       uintptr(0),
							File:     "",
							Line:     0,
							Function: "",
						},
						Stack: "",
					},
					Context: []zapcore.Field{
						{
							Key:       "component",
							Type:      zapcore.StringType,
							Integer:   0,
							String:    "dummyscheduler",
							Interface: nil,
						},
					},
				},
				{
					Entry: zapcore.Entry{
						Level:      zapcore.ErrorLevel,
						Time:       time.Time{},
						LoggerName: "",
						Message:    "what happened",
						Caller: zapcore.EntryCaller{
							Defined:  false,
							PC:       uintptr(0),
							File:     "",
							Line:     0,
							Function: "",
						},
						Stack: "",
					},
					Context: []zapcore.Field{
						{
							Key:       "component",
							Type:      zapcore.StringType,
							Integer:   0,
							String:    "dummyscheduler",
							Interface: nil,
						},
					},
				},
			},
			testFunc: func(t *testing.T, l *logrus.Logger, r *observer.ObservedLogs) {
				l.WithField("component", "dummyscheduler").Info("here you go")
				l.WithField("component", "dummyscheduler").Error("what happened")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logrusLogger.SetOutput(ioutil.Discard)

			core, recorded := observer.New(zapcore.InfoLevel)
			logger := zap.New(core)

			hook, _ := NewLogrusZapHook(logger)
			tt.logrusLogger.AddHook(hook)

			tt.testFunc(t, tt.logrusLogger, recorded)

			fixedTime := time.Now()

			trans := cmp.Transformer("date", func(in []observer.LoggedEntry) []observer.LoggedEntry {
				out := []observer.LoggedEntry{}
				for _, e := range in {
					e.Time = fixedTime
					out = append(out, e)
				}
				return out
			})

			if !cmp.Equal(tt.expected, recorded.All(), trans) {
				t.Errorf("expected log entries and recorded entries do not match: %s", cmp.Diff(tt.expected, recorded.All(), trans))
			}

		})
	}
}
