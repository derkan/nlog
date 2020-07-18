package json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/formatter"
	"github.com/derkan/nlog/loader"
	"github.com/derkan/nlog/pool"
	"github.com/derkan/nlog/writer"
	fl "github.com/derkan/nlog/writer/filerotater"
)

var (
	// TimeKey is for key for time
	TimeKey = "time"
	// PrefixKey is key name for logger logger name
	NameKey = "logger"
	// LevelKey is key name for logging level
	LevelKey = "level"
	// LocKey is for file location info key
	LocKey = "loc"
	// MsgKey is for key for message
	MsgKey = "msg"
)

// Formatter logs with json format
type Formatter struct {
	cfg   *config
	hooks []nlog.Hook
}

// NewFormatter returns a new instance of Formatter
func NewFormatter(opts ...option) *Formatter {
	c := &Formatter{
		cfg: &config{Level: nlog.INFO, FileLocCallerDepth: 4}}
	// Loop through each option and set
	for _, opt := range opts {
		opt(c.cfg)
	}
	c.SetDefaults()
	return c
}

// NewFromConfig builds a JSON formatter from given loader config
// appName is used in syslog
func NewFromConfig(f loader.Formatter, appName string) *Formatter {
	c := &Formatter{cfg: &config{
		Level:              f.Level,
		Date:               f.Date,
		FileLoc:            f.FileLoc,
		FileLocStrip:       f.FileLocStrip,
		FileLocCallerDepth: f.FileLocCallerDepth,
		NoPrintLevel:       f.NoPrintLevel,
		Time:               f.Time,
		TimeResolution:     f.TimeResolution,
		TimeUTC:            f.TimeUTC,
		UnixTime:           f.UnixTime,
	}}

	var parallelW []*writer.ParallelWriter
	var normalW []*writer.Writer
	var wrt io.WriteCloser
	for _, w := range f.Writers {
		switch w.Type {
		case "stdout":
			wrt = os.Stdout
		case "stderr":
			wrt = os.Stderr
		case "syslog":
			wrt = writer.NewSysLogWriter(appName)
		case "filerotator":
			wrt = &fl.Rotater{
				Filename:   w.Filename,
				Compress:   w.Compress,
				MaxAge:     w.MaxAge,
				LocalTime:  !w.UTC,
				MaxBackups: w.MaxBackups,
				MaxSize:    w.MaxSize,
			}
		default:
			continue
		}
		if f.LeveledType == "parallel" {
			parallelW = append(parallelW, writer.NewParellelWriter(wrt, w.Level, w.QueueLen))
		} else {
			normalW = append(normalW, writer.NewWriter(wrt, w.Level))
		}
	}
	if len(normalW) > 0 {
		c.cfg.Writer = writer.NewMultiWriter(normalW...)
	}
	if len(parallelW) > 0 {
		c.cfg.Writer = writer.NewParallelMultiWriter(parallelW...)
	}
	c.SetDefaults()
	return c
}

func (cl *Formatter) SetDefaults() {
	if cl.cfg.Writer == nil {
		cl.cfg.Writer = writer.NewMultiWriter(writer.NewWriter(os.Stderr, cl.cfg.Level))
	}
	if cl.cfg.MarshallFn == nil {
		cl.cfg.MarshallFn = json.Marshal
	}
}
func (cl *Formatter) GetCallDepth(sub int) int {
	return cl.cfg.FileLocCallerDepth - sub
}

// Init inits formatter
func (cl *Formatter) Init() {
}

// Flush flushes to disk and closes writers
func (cl *Formatter) Flush() {
	cl.cfg.Writer.Close()
}

// levelStr gets formatted level string
func (cl *Formatter) levelStr(lvl nlog.Level) string {
	if v, ok := nlog.LevelNames[lvl]; ok {
		return v
	}
	return fmt.Sprintf("!%d", lvl)
}

func (cl *Formatter) hookSet(buff nlog.Buffer, lvl nlog.Level) nlog.HookFieldFn {
	return func(key string, value interface{}) {
		cl.AppendKV(buff, lvl, key, value)
	}
}

// Logf logs current log line without with args format
func (cl *Formatter) Logf(callDepth int, lvl nlog.Level, fields nlog.Buffer, lgName, layout string, args ...interface{}) {
	if callDepth == 0 {
		callDepth = cl.cfg.FileLocCallerDepth
	}
	buff := pool.GetBuffer()
	defer pool.PutBuffer(buff)
	// time
	buff.AppendByte('{')
	cl.cfg.GetTime(buff, true)
	if buff.Len() > 1 {
		buff.AppendByte(',')
	}
	// level
	if !cl.cfg.NoPrintLevel {
		cl.AppendFieldKey(buff, LevelKey)
		cl.AppendFieldValue(buff, cl.levelStr(lvl))
	}
	// Logger name
	if lgName != "" {
		buff.AppendByte(',')
		cl.AppendFieldKey(buff, NameKey)
		cl.AppendFieldValue(buff, lgName)
	}

	// Message
	buff.AppendByte(',')
	cl.AppendFieldKey(buff, MsgKey)
	if layout != "" {
		fmt.Fprintf(buff, "%q", fmt.Sprintf(layout, args...))
	} else {
		fmt.Fprintf(buff, "%q", fmt.Sprint(args...))
	}

	// Let assigning new fields from hooks even if fields is nil
	var hookBuff nlog.Buffer
	if len(cl.cfg.Hooks) > 0 {
		var buffSet nlog.HookBufferSet
		if fields == nil {
			hookBuff = pool.GetBuffer()
			buffSet = nlog.HookBufferSet{
				Buffer: hookBuff,
				With:   cl.hookSet(hookBuff, lvl),
			}
		} else {
			buffSet = nlog.HookBufferSet{
				Buffer: fields,
				With:   cl.hookSet(fields, lvl),
			}
		}
		// Run hooks
		for i := range cl.cfg.Hooks {
			cl.cfg.Hooks[i].Run(lvl, buffSet, fmt.Sprintf(layout, args...))
		}
	}
	// Fields filled only from hooks
	if hookBuff != nil {
		hookBuff.WriteTo(buff)
		pool.PutBuffer(hookBuff)
	}
	// Fields
	if fields != nil {
		fields.WriteTo(buff)
	}
	// File location
	if cl.cfg.FileLoc {
		buff.AppendByte(',')
		cl.AppendFieldKey(buff, LocKey)
		formatter.GetFileLoc(cl.cfg.FileLocStrip, buff, callDepth, true)
	}
	buff.AppendByte('}')
	buff.AppendByte('\n')
	cl.cfg.Writer.WriteIfLevel(lvl, buff.Bytes())

}

// AppendFieldKey appends key value to log
func (cl *Formatter) AppendFieldKey(buf nlog.Buffer, key string) {
	fmt.Fprintf(buf, "%q:", key)
}

// AppendFieldValue appends value to log
func (cl *Formatter) AppendFieldValue(buf nlog.Buffer, val interface{}) {
	buf.AppendAny(val, true, cl.cfg.MarshallFn)
}

// AppendKV formats and appends key,value to buffer
func (cl *Formatter) AppendKV(buff nlog.Buffer, lvl nlog.Level, key string, val interface{}) {
	if cl.cfg.Level < lvl {
		return
	}
	buff.AppendByte(',')
	cl.AppendFieldKey(buff, key)
	buff.AppendAny(val, true, cl.cfg.MarshallFn)
}

// Str appends str value to buff with a format
func (cl *Formatter) Str(buf nlog.Buffer, lvl nlog.Level, key, val string) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	fmt.Fprintf(buf, "%q", val)
}

// Strs adds a slice of string value with a key to buff
func (cl *Formatter) Strs(buf nlog.Buffer, lvl nlog.Level, key string, val []string) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendStrings(val, true)
}

// Int adds a new int key value to buff
func (cl *Formatter) Int(buf nlog.Buffer, lvl nlog.Level, key string, val int) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInt(val)
}

// Ints adds a slice of int value with a key to buff
func (cl *Formatter) Ints(buf nlog.Buffer, lvl nlog.Level, key string, val []int) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInts(val)
}

// Ints8 adds a slice of int8 value with a key to buff
func (cl *Formatter) Ints8(buf nlog.Buffer, lvl nlog.Level, key string, val []int8) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInts8(val)
}

// Ints16 adds a slice of int16 value with a key to buff
func (cl *Formatter) Ints16(buf nlog.Buffer, lvl nlog.Level, key string, val []int16) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInts16(val)
}

// Ints32 adds a slice of int32 value with a key to buff
func (cl *Formatter) Ints32(buf nlog.Buffer, lvl nlog.Level, key string, val []int32) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInts32(val)
}

// Int64 adds a new int64 key value to buff
func (cl *Formatter) Int64(buf nlog.Buffer, lvl nlog.Level, key string, val int64) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInt64(val)
}

// Int64s adds a slice of int64 value with a key to buff
func (cl *Formatter) Int64s(buf nlog.Buffer, lvl nlog.Level, key string, val []int64) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendInts64(val)
}

// UInt adds a new uint key value to buff
func (cl *Formatter) UInt(buf nlog.Buffer, lvl nlog.Level, key string, val uint) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendUInt64(uint64(val))
}

// UInts adds a slice of uint value with a key to buff
func (cl *Formatter) UInts(buf nlog.Buffer, lvl nlog.Level, key string, val []uint) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendUInts(val)
}

// UInts16 adds a slice of uint16 value with a key to buff
func (cl *Formatter) UInts16(buf nlog.Buffer, lvl nlog.Level, key string, val []uint16) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendUInts16(val)
}

// UInts32 adds a slice of uint32 value with a key to buff
func (cl *Formatter) UInts32(buf nlog.Buffer, lvl nlog.Level, key string, val []uint32) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendUInts32(val)
}

// Uint64 adds a new uint64 key value to buff
func (cl *Formatter) UInt64(buf nlog.Buffer, lvl nlog.Level, key string, val uint64) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendUInt64(val)
}

// UInts64 adds a slice of uint64 value with a key to buff
func (cl *Formatter) UInts64(buf nlog.Buffer, lvl nlog.Level, key string, val []uint64) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendUInts64(val)
}

// Float32 adds a new float32 key value to buff
func (cl *Formatter) Float32(buf nlog.Buffer, lvl nlog.Level, key string, val float32) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendFloat32(val)
}

// Floats32 adds a slice of float32 value with a key to buff
func (cl *Formatter) Floats32(buf nlog.Buffer, lvl nlog.Level, key string, val []float32) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendFloats32(val)
}

// Float64 adds a new float64 key value to buff
func (cl *Formatter) Float64(buf nlog.Buffer, lvl nlog.Level, key string, val float64) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendFloat64(val)
}

// Floats64 adds a slice of float32 value with a key to buff
func (cl *Formatter) Floats64(buf nlog.Buffer, lvl nlog.Level, key string, val []float64) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendFloats64(val)
}

// Bool adds a new bool key value to buff
func (cl *Formatter) Bool(buf nlog.Buffer, lvl nlog.Level, key string, val bool) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendBool(val)
}

// Bools adds a slice of bool value with a key to buff
func (cl *Formatter) Bools(buf nlog.Buffer, lvl nlog.Level, key string, val []bool) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendBools(val)
}

// Error adds a new error key value to buff
func (cl *Formatter) Error(buf nlog.Buffer, lvl nlog.Level, key string, val error) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendError(val, true)
}

// Bools adds a slice of error value with a key to buff
func (cl *Formatter) Errors(buf nlog.Buffer, lvl nlog.Level, key string, val []error) {
	if cl.cfg.Level < lvl {
		return
	}
	buf.AppendByte(',')
	fmt.Fprintf(buf, "%q:", key)
	buf.AppendErrors(val, true)
}
