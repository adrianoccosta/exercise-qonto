package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	zapLog *zap.Logger
}

func New() (Logger, error) {
	var err error
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	enccoderConfig.StacktraceKey = "" // to hide stacktrace info
	config.EncoderConfig = enccoderConfig

	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &logger{
		zapLog: log,
	}, nil
}

// Logger defines the common interface for the logs
type Logger interface {

	// With returns the current logger with the provided fields as metadata to be logged.
	With(fields ...zap.Field) Logger

	// WithError returns the current logger with the provided error as metadata to be logged.
	WithError(err error) Logger

	// Debug prints the msg in debug level, with the provided fields metadata to be logged.
	Debug(msg string, fields ...zap.Field)

	// Info prints the msg in info level, with the provided fields metadata to be logged.
	Info(msg string, fields ...zap.Field)

	// Warn prints the msg in warn level, with the provided fields metadata to be logged.
	Warn(msg string, fields ...zap.Field)

	// Error prints the msg in error level, with the provided fields metadata to be logged.
	Error(msg string, fields ...zap.Field)

	// Fatal prints the msg in fatal level, with the provided fields metadata to be logged.
	Fatal(msg string, fields ...zap.Field)

	// Panic prints the msg in panic level, with the provided fields metadata to be logged.
	Panic(msg string, fields ...zap.Field)
}

func (l logger) With(fields ...zap.Field) Logger {
	return &logger{
		zapLog: l.zapLog.With(fields...),
	}
}

func (l logger) WithError(err error) Logger {
	return &logger{
		zapLog: l.zapLog.With(zap.Any("error", err)),
	}
}

func (l logger) Debug(message string, fields ...zap.Field) {
	l.zapLog.Debug(message, fields...)
}

func (l logger) Info(message string, fields ...zap.Field) {
	l.zapLog.Info(message, fields...)
}

func (l logger) Warn(message string, fields ...zap.Field) {
	l.zapLog.Warn(message, fields...)
}

func (l logger) Error(message string, fields ...zap.Field) {
	l.zapLog.Error(message, fields...)
}

func (l logger) Fatal(message string, fields ...zap.Field) {
	l.zapLog.Fatal(message, fields...)
}

func (l logger) Panic(message string, fields ...zap.Field) {
	l.zapLog.Panic(message, fields...)
}
