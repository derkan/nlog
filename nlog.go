package nlog

import (
	"io"
)

// Logger represent common interface for logging function
type Logger interface {
	// Flush closes all writers safely
	Flush()
	// Fatal returns FATAL level logger item
	Fatal() LoggerItem
	// Fatalf logs FATAL level log with given format-params and exits
	Fatalf(format string, args ...interface{})
	// Error returns ERROR level logger item
	Error() LoggerItem
	// Errorf logs ERROR level log with given format-params
	Errorf(format string, args ...interface{})
	// Warn returns WARNING level logger item
	Warn() LoggerItem
	// Warnf logs WARN level log with given format-params
	Warnf(format string, args ...interface{})
	// Info returns info level logger item
	Info() LoggerItem
	// Infof prints INFO level message with given format and args
	Infof(format string, args ...interface{})
	// Debug returns DEBUG level logger item
	Debug() LoggerItem
	// Debugf prints DEBUG level message with given format-params
	Debugf(format string, args ...interface{})
	// Print prints log at INFO level with given params
	Print(args ...interface{})
	// Printf prints log at INFO level with given format-params
	Printf(format string, args ...interface{})
	// Print prints log at INFO level with given params
	Println(args ...interface{})
	// Sub returns a sub logger with given prefix and optionally min logging level
	Sub(prefix string, minLevel ...Level) Logger
}

// MarshallFn is func stub used for custom marshalling
type MarshallFn func(interface{}) ([]byte, error)

// LoggerItem is logger context
type LoggerItem interface {
	// Msg logs info message with format if logging level is satisfied
	Msg(args ...interface{})
	// Msgf logs info message with format if logging level is satisfied
	Msgf(format string, args ...interface{})
	// Str adds a new str key value to buff, Msg/Msgf should be called in same chain
	Str(key string, val string) LoggerItem
	// Strs adds a slice of string value with a key to buff, Msg/Msgf should be called in same chain
	Strs(key string, val []string) LoggerItem
	// Int adds a new int key value to buff, Msg/Msgf should be called in same chain
	Int(key string, val int) LoggerItem
	// Ints adds a slice of int value with a key to buff, Msg/Msgf should be called in same chain
	Ints(key string, val []int) LoggerItem
	// Ints8 adds a slice of int8 value with a key to buff, Msg/Msgf should be called in same chain
	Ints8(key string, val []int8) LoggerItem
	// Ints16 adds a slice of int16 value with a key to buff, Msg/Msgf should be called in same chain
	Ints16(key string, val []int16) LoggerItem
	// Ints32 adds a slice of int32 value with a key to buff, Msg/Msgf should be called in same chain
	Ints32(key string, val []int32) LoggerItem
	// Int64 adds a new int64 key value to buff, Msg/Msgf should be called in same chain
	Int64(key string, val int64) LoggerItem
	// Int64s adds a slice of int64 value with a key to buff, Msg/Msgf should be called in same chain
	Int64s(key string, val []int64) LoggerItem
	// UInt adds a new uint key value to buff, Msg/Msgf should be called in same chain
	UInt(key string, val uint) LoggerItem
	// UInts adds a slice of uint value with a key to buff, Msg/Msgf should be called in same chain
	UInts(key string, val []uint) LoggerItem
	// UInts16 adds a slice of uint16 value with a key to buff, Msg/Msgf should be called in same chain
	UInts16(key string, val []uint16) LoggerItem
	// UInts32 adds a slice of uint32 value with a key to buff, Msg/Msgf should be called in same chain
	UInts32(key string, val []uint32) LoggerItem
	// UInt64 adds a new uint64 key value to buff, Msg/Msgf should be called in same chain
	UInt64(key string, val uint64) LoggerItem
	// UInts64 adds a slice of uint64 value with a key to buff, Msg/Msgf should be called in same chain
	UInts64(key string, val []uint64) LoggerItem
	// Float32 adds a new float32 key value to buff, Msg/Msgf should be called in same chain
	Float32(key string, val float32) LoggerItem
	// Floats32 adds a slice of float32 value with a key to buff, Msg/Msgf should be called in same chain
	Floats32(key string, val []float32) LoggerItem
	// Float64 adds a new float64 key value to buff, Msg/Msgf should be called in same chain
	Float64(key string, val float64) LoggerItem
	// Floats64 adds a slice of float32 value with a key to buff, Msg/Msgf should be called in same chain
	Floats64(key string, val []float64) LoggerItem
	// Bool adds a new bool key value to buff, Msg/Msgf should be called in same chain
	Bool(key string, val bool) LoggerItem
	// Bools adds a slice of bool value with a key to buff, Msg/Msgf should be called in same chain
	Bools(key string, val []bool) LoggerItem
	// Error adds a new error key value to buff, Msg/Msgf should be called in same chain
	Err(val error) LoggerItem
	// Bools adds a slice of error value with a key to buff, Msg/Msgf should be called in same chain
	Errors(key string, val []error) LoggerItem
	// With adds a new str key value to buff, Msg/Msgf should be called in same chain
	With(key string, val interface{}) LoggerItem
}

// Buffer provides byte buffer, which can be used for minimizing
// memory allocations.
type Buffer interface {
	// Cap returns the capacity of the byte buffer.
	Cap() int
	// Len returns the size of the byte buffer.
	Len() int
	// The function appends all the data read from r to b.
	ReadFrom(r io.Reader) (int64, error)
	// WriteTo implements io.WriterTo.
	WriteTo(w io.Writer) (int64, error)
	// Bytes returns buffer Bytes , i.e. all the bytes accumulated in the buffer.
	Bytes() []byte
	// Write implements io.Writer - it appends p to Buffer.B
	Write(p []byte) (int, error)
	// Set sets Buffer Bytes to p.
	Set(p []byte)
	// SetString sets Buffer bytes to s.
	SetString(s string)
	// String returns string representation of Buffer.B.
	String() string
	// Reset makes Buffer empty.
	Reset()
	// Itoa is Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
	Itoa(i int, wid int)
	// AppendByte writes a single byte to the Buffer.
	AppendByte(val byte) Buffer
	// AppendBytes writes a slice of byte to the Buffer.
	AppendBytes(vals []byte) Buffer
	// AppendString writes a string to the Buffer.
	AppendString(val string, quota bool) Buffer
	// AppendStrings writes a slice of string to the Buffer.
	AppendStrings(vals []string, quota bool) Buffer
	// AppendInt appends an integer to the underlying buffer (assuming base 10).
	AppendInt(val int) Buffer
	// AppendInts appends a slice of integer to the underlying buffer
	AppendUInt64(val uint64) Buffer
	// AppendInt64 appends an int64 to the underlying buffer (assuming base 10).
	AppendInt64(val int64) Buffer
	// AppendUInts64 appends a slice of uint64 to the underlying buffer (assuming base 10).
	AppendInts(vals []int) Buffer
	// AppendInts8 appends a slice of int8 to the underlying buffer
	AppendInts8(val []int8) Buffer
	// AppendInts16 appends a slice of int16 to the underlying buffer
	AppendInts16(val []int16) Buffer
	// AppendInts32 appends a slice of int32 to the underlying buffer
	AppendInts32(val []int32) Buffer
	// AppendInt64 appends a slice of int64 to the underlying buffer
	AppendInts64(val []int64) Buffer
	// AppendUInts appends a slice of integer to the underlying buffer
	AppendUInts(vals []uint) Buffer
	// AppendUInts16 appends a slice of int16 to the underlying buffer
	AppendUInts16(val []uint16) Buffer
	// AppendUInts32 appends a slice of int32 to the underlying buffer
	AppendUInts32(val []uint32) Buffer
	// AppendUInt64 appends an uint64 to the underlying buffer (assuming base 10).
	AppendUInts64(vals []uint64) Buffer
	// AppendFloat32 appends an float32 to the underlying buffer
	AppendFloat32(val float32) Buffer
	// AppendFloats32 appends a slice of float32 to the underlying buffer
	AppendFloats32(vals []float32) Buffer
	// AppendFloat64 appends an float64 to the underlying buffer
	AppendFloat64(val float64) Buffer
	// AppendFloats64 appends a slice of float64 to the underlying buffer
	AppendFloats64(vals []float64) Buffer
	// AppendBool appends a bool to the underlying buffer.
	AppendBool(val bool) Buffer
	// AppendBools appends a slice bool to the underlying buffer.
	AppendBools(vals []bool) Buffer
	// AppendError writes error to buffer
	AppendError(val error, quota bool) Buffer
	// AppendErrors writes a slice of string to the Buffer.
	AppendErrors(vals []error, quota bool) Buffer
	// AppendInterface takes an arbitrary object and converts it to JSON and embeds it dst.
	AppendInterface(val interface{}, marshallFn MarshallFn) Buffer
	// AppendAny appends given interface with its type
	AppendAny(val interface{}, quota bool, marshallFn MarshallFn) Buffer
}

type HookFieldFn func(key string, value interface{})

type HookBufferSet struct {
	Buffer Buffer
	With   HookFieldFn
}

// Hook defines an interface to a log hook.
type Hook interface {
	// Run runs the hook with the event.
	Run(level Level, buffer HookBufferSet, message string)
}

// HookFunc is an adaptor to allow the use of an ordinary function
// as a Hook.
type HookFunc func(level Level, buffer HookBufferSet, message string)

// Run implements the Hook interface.
func (h HookFunc) Run(level Level, buffer HookBufferSet, message string) {
	h(level, buffer, message)
}

// Formatter is interface type for defining log output formatters
type Formatter interface {
	// Init inits formatter with given config
	Init()
	// Flush flushes to disk and closes writers
	Flush()
	// Logf logs current log line with args format
	Logf(callDepth int, lvl Level, buff Buffer, lgName, format string, args ...interface{})
	// AddField appends key=value
	AppendKV(buf Buffer, lvl Level, key string, val interface{})
	// Str appends str value to buff with a format
	Str(buf Buffer, lvl Level, key string, val string)
	// Strs adds a slice of string value with a key to buff
	Strs(buf Buffer, lvl Level, key string, val []string)
	// Int adds a new int key value to buff
	Int(buf Buffer, lvl Level, key string, val int)
	// Ints adds a slice of int value with a key to buff
	Ints(buf Buffer, lvl Level, key string, val []int)
	// Ints8 adds a slice of int8 value with a key to buff
	Ints8(buf Buffer, lvl Level, key string, val []int8)
	// Ints16 adds a slice of int16 value with a key to buff
	Ints16(buf Buffer, lvl Level, key string, val []int16)
	// Ints32 adds a slice of int32 value with a key to buff
	Ints32(buf Buffer, lvl Level, key string, val []int32)
	// Int64 adds a new int64 key value to buff
	Int64(buf Buffer, lvl Level, key string, val int64)
	// Int64s adds a slice of int64 value with a key to buff
	Int64s(buf Buffer, lvl Level, key string, val []int64)
	// UInt adds a new uint key value to buff
	UInt(buf Buffer, lvl Level, key string, val uint)
	// UInts adds a slice of uint value with a key to buff
	UInts(buf Buffer, lvl Level, key string, val []uint)
	// UInts16 adds a slice of uint16 value with a key to buff
	UInts16(buf Buffer, lvl Level, key string, val []uint16)
	// UInts32 adds a slice of uint32 value with a key to buff
	UInts32(buf Buffer, lvl Level, key string, val []uint32)
	// UInt64 adds a new uint64 key value to buff
	UInt64(buf Buffer, lvl Level, key string, val uint64)
	// UInts64 adds a slice of uint64 value with a key to buff
	UInts64(buf Buffer, lvl Level, key string, val []uint64)
	// Float32 adds a new float32 key value to buff
	Float32(buf Buffer, lvl Level, key string, val float32)
	// Floats32 adds a slice of float32 value with a key to buff
	Floats32(buf Buffer, lvl Level, key string, val []float32)
	// Float64 adds a new float64 key value to buff
	Float64(buf Buffer, lvl Level, key string, val float64)
	// Floats64 adds a slice of float32 value with a key to buff
	Floats64(buf Buffer, lvl Level, key string, val []float64)
	// Bool adds a new bool key value to buff
	Bool(buf Buffer, lvl Level, key string, val bool)
	// Bools adds a slice of bool value with a key to buff
	Bools(buf Buffer, lvl Level, key string, val []bool)
	// Error adds a new error key value to buff
	Error(buf Buffer, lvl Level, key string, val error)
	// Bools adds a slice of error value with a key to buff
	Errors(buf Buffer, lvl Level, key string, val []error)
	//  GetCallDepth returns call depth
	GetCallDepth(sub int) int
}
