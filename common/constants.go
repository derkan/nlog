package common

// Level defines all available log levels for log messages.
type Level int

// Log levels.
const (
	FATAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var (
	// DebugStr is string value for debug log level key
	DebugStr = "DBG"
	// InfoStr is string value for info log level key
	InfoStr = "INF"
	// WarnStr is string value for warn log level key
	WarnStr = "WRN"
	// ErrorStr is string value for error log level key
	ErrorStr = "ERR"
	// FatalStr is string value for falal log level key
	FatalStr = "FAT"
)

var LevelNames = map[Level]string{
	FATAL:   FatalStr,
	ERROR:   ErrorStr,
	WARNING: WarnStr,
	INFO:    InfoStr,
	DEBUG:   DebugStr,
}
