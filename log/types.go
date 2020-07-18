package log

import "github.com/derkan/nlog"

// Log is current logger instance
var Log Logger

// Logger represent common interface for logging function
type Logger interface {
	// Flush closes all writers safely
	Flush()
	// Fatal returns FATAL level logger item
	Fatal() nlog.LoggerItem
	// Fatalf logs FATAL level log with given format-params and exits
	Fatalf(format string, args ...interface{})
	// Error returns ERROR level logger item
	Error() nlog.LoggerItem
	// Errorf logs ERROR level log with given format-params
	Errorf(format string, args ...interface{})
	// Warn returns WARNING level logger item
	Warn() nlog.LoggerItem
	// Warnf logs WARN level log with given format-params
	Warnf(format string, args ...interface{})
	// Info returns info level logger item
	Info() nlog.LoggerItem
	// Infof prints INFO level message with given format and args
	Infof(format string, args ...interface{})
	// Debug returns DEBUG level logger item
	Debug() nlog.LoggerItem
	// Debugf prints DEBUG level message with given format-params
	Debugf(format string, args ...interface{})
	// Print prints log at INFO level with given params
	Print(args ...interface{})
	// Printf prints log at INFO level with given format-params
	Printf(format string, args ...interface{})
	// Print prints log at INFO level with given params
	Println(args ...interface{})
	// Sub returns a sub logger with given prefix and optionally min logging level
	Sub(prefix string, minLevel ...nlog.Level) Logger
}
