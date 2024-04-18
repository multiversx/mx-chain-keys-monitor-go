package mock

import (
	logger "github.com/multiversx/mx-chain-logger-go"
)

// LoggerStub -
type LoggerStub struct {
	TraceHandler      func(message string, args ...interface{})
	DebugHandler      func(message string, args ...interface{})
	InfoHandler       func(message string, args ...interface{})
	WarnHandler       func(message string, args ...interface{})
	ErrorHandler      func(message string, args ...interface{})
	LogIfErrorHandler func(err error, args ...interface{})
	LogHandler        func(logLevel logger.LogLevel, message string, args ...interface{})
	LogLineHandler    func(line *logger.LogLine)
	SetLevelHandler   func(logLevel logger.LogLevel)
	GetLevelHandler   func() logger.LogLevel
}

// Trace -
func (stub *LoggerStub) Trace(message string, args ...interface{}) {
	if stub.TraceHandler != nil {
		stub.TraceHandler(message, args...)
	}
}

// Debug -
func (stub *LoggerStub) Debug(message string, args ...interface{}) {
	if stub.DebugHandler != nil {
		stub.DebugHandler(message, args...)
	}
}

// Info -
func (stub *LoggerStub) Info(message string, args ...interface{}) {
	if stub.InfoHandler != nil {
		stub.InfoHandler(message, args...)
	}
}

// Warn -
func (stub *LoggerStub) Warn(message string, args ...interface{}) {
	if stub.WarnHandler != nil {
		stub.WarnHandler(message, args...)
	}
}

// Error -
func (stub *LoggerStub) Error(message string, args ...interface{}) {
	if stub.ErrorHandler != nil {
		stub.ErrorHandler(message, args...)
	}
}

// LogIfError -
func (stub *LoggerStub) LogIfError(err error, args ...interface{}) {
	if stub.LogIfErrorHandler != nil {
		stub.LogIfErrorHandler(err, args...)
	}
}

// Log -
func (stub *LoggerStub) Log(logLevel logger.LogLevel, message string, args ...interface{}) {
	if stub.LogHandler != nil {
		stub.LogHandler(logLevel, message, args...)
	}
}

// LogLine -
func (stub *LoggerStub) LogLine(line *logger.LogLine) {
	if stub.LogLineHandler != nil {
		stub.LogLineHandler(line)
	}
}

// SetLevel -
func (stub *LoggerStub) SetLevel(logLevel logger.LogLevel) {
	if stub.SetLevelHandler != nil {
		stub.SetLevelHandler(logLevel)
	}
}

// GetLevel -
func (stub *LoggerStub) GetLevel() logger.LogLevel {
	if stub.GetLevelHandler != nil {
		return stub.GetLevelHandler()
	}

	return logger.LogNone
}

// IsInterfaceNil -
func (stub *LoggerStub) IsInterfaceNil() bool {
	return stub == nil
}
