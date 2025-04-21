package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	// Set output to stdout
	Log.SetOutput(os.Stdout)

	// Set log level
	Log.SetLevel(logrus.InfoLevel)

	// Set formatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		DisableColors:   false,
	})
}

// WithFields creates a new entry with fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Fatal logs a message at level Fatal
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Panic logs a message at level Panic
func Panic(args ...interface{}) {
	Log.Panic(args...)
}
