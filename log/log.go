package log

import "github.com/derkan/nlog"

// Logger is default logger instance
var Logger nlog.Logger

func init() {
	Logger = New()
}

// Flush closes all writers safely
func Flush() {
	Logger.Flush()
}

// Fatal returns FATAL level logger item
func Fatal() nlog.LoggerItem {
	return Logger.Fatal()
}

// Fatalf logs FATAL level log with given format-params and exits
func Fatalf(msg string, v ...interface{}) {
	Logger.Fatalf(msg, v...)
}

// Error returns ERROR level logger item
func Error() nlog.LoggerItem {
	return Logger.Error()
}

// Errorf logs ERROR level log with given format-params
func Errorf(msg string, v ...interface{}) {
	Logger.Errorf(msg, v...)
}

// Warn returns WARNING level logger item
func Warn() nlog.LoggerItem {
	return Logger.Warn()
}

// Warnf logs WARN level log with given format-params
func Warnf(msg string, v ...interface{}) {
	Logger.Warnf(msg, v...)
}

// Info returns info level logger item
func Info() nlog.LoggerItem {
	return Logger.Info()
}

// Infof prints INFO level message with given format and args
func Infof(msg string, v ...interface{}) {
	Logger.Infof(msg, v...)
}

// Debug returns DEBUG level logger item
func Debug() nlog.LoggerItem {
	return Logger.Debug()
}

// Debugf prints DEBUG level message with given format-params
func Debugf(msg string, v ...interface{}) {
	Logger.Debugf(msg, v...)
}

// Print prints log at INFO level with given params
func Print(msg ...interface{}) {
	Logger.Print(msg...)
}

// Printf prints log at INFO level with given format-params
func Printf(msg string, v ...interface{}) {
	Logger.Printf(msg, v...)
}

// Print prints log at INFO level with given params
func Println(msg ...interface{}) {
	Logger.Print(msg...)
}

// Sub returns a sub logger with given prefix and optionally min logging level
func Sub(prefix string, minLevel ...nlog.Level) nlog.Logger {
	return Logger.Sub(prefix, minLevel...)
}
