package log

import "github.com/derkan/nlog/common"

// StdLogger is standard logger instance
var StdLogger Logger

func init() {
	StdLogger = New()
}

// Flush closes all writers safely
func Flush() {
	StdLogger.Flush()
}

// Fatal returns FATAL level logger item
func Fatal() common.LoggerItem {
	return StdLogger.Fatal()
}

// Fatalf logs FATAL level log with given format-params and exits
func Fatalf(msg string, v ...interface{}) {
	StdLogger.Fatalf(msg, v...)
}

// Error returns ERROR level logger item
func Error() common.LoggerItem {
	return StdLogger.Error()
}

// Errorf logs ERROR level log with given format-params
func Errorf(msg string, v ...interface{}) {
	StdLogger.Errorf(msg, v...)
}

// Warn returns WARNING level logger item
func Warn() common.LoggerItem {
	return StdLogger.Warn()
}

// Warnf logs WARN level log with given format-params
func Warnf(msg string, v ...interface{}) {
	StdLogger.Warnf(msg, v...)
}

// Info returns info level logger item
func Info() common.LoggerItem {
	return StdLogger.Info()
}

// Infof prints INFO level message with given format and args
func Infof(msg string, v ...interface{}) {
	StdLogger.Infof(msg, v...)
}

// Debug returns DEBUG level logger item
func Debug() common.LoggerItem {
	return StdLogger.Debug()
}

// Debugf prints DEBUG level message with given format-params
func Debugf(msg string, v ...interface{}) {
	StdLogger.Debugf(msg, v...)
}

// Print prints log at INFO level with given params
func Print(msg ...interface{}) {
	StdLogger.Print(msg...)
}

// Printf prints log at INFO level with given format-params
func Printf(msg string, v ...interface{}) {
	StdLogger.Printf(msg, v...)
}

// Print prints log at INFO level with given params
func Println(msg ...interface{}) {
	StdLogger.Print(msg...)
}

// Sub returns a sub logger with given prefix and optionally min logging level
func Sub(prefix string, minLevel ...common.Level) Logger {
	return StdLogger.Sub(prefix, minLevel...)
}
