// Package logger implements general logger.
package logger

//go:generate go run github.com/golang/mock/mockgen --source=logger.go --destination=logger_mock.go --package=logger

// Logger describes general logging methods.
type Logger interface {
	LogWriter
}

// LogWriter describes final log-write methods.
type LogWriter interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}
