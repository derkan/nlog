package console

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

// Formatter logs pretty
type Formatter struct {
	cfg       *config
	debugStr  string
	infoStr   string
	noticeStr string
	warnStr   string
	errorStr  string
	fatalStr  string
	strW      func(nlog.Buffer, string, string)
	intW      func(nlog.Buffer, string, int64)
	hooks     []nlog.Hook
}

// NewFormatter returns a new instance of Formatter
func NewFormatter(opts ...option) *Formatter {
	c := &Formatter{cfg: &config{Level: nlog.INFO, FileLocCallerDepth: 4}}
	// Loop through each option and set
	for _, opt := range opts {
		opt(c.cfg)
	}
	c.SetDefaults()
	return c
}

// NewFromConfig builds a console formatter from given loader config
// appName is used in syslog
func NewFromConfig(f loader.Formatter, appName string) *Formatter {
	c := &Formatter{cfg: &config{
		Colored:            f.Colored,
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
	if cl.cfg.FileLocCallerDepth < 3 {
		cl.cfg.FileLocCallerDepth = 4
	}
}
func (cl *Formatter) GetCallDepth(sub int) int {
	return cl.cfg.FileLocCallerDepth - sub
}

// Flush flushes to disk and closes writers
func (cl *Formatter) Flush() {
	cl.cfg.Writer.Close()
}

func strC(b nlog.Buffer, c string, v string) {
	b.AppendString(c, false).AppendString(v, false).AppendString(ColorReset, false)
}

func strW(b nlog.Buffer, c string, v string) {
	b.AppendString(v, false)
}

func intC(b nlog.Buffer, c string, v int64) {
	b.AppendString(c, false).AppendInt64(v).AppendString(ColorReset, false)
}

func intW(b nlog.Buffer, c string, v int64) {
	b.AppendInt64(v)
}

// Init initializes formatter
func (cl *Formatter) Init() {
	if cl.cfg.Colored {
		initColor(cl.cfg.Writer)
		cl.strW = strC
		cl.intW = intC
		cl.debugStr = fmt.Sprintf("%s%s%s ", DebugColor, nlog.DebugStr, ColorReset)
		cl.infoStr = fmt.Sprintf("%s%s%s ", InfoColor, nlog.InfoStr, ColorReset)
		cl.warnStr = fmt.Sprintf("%s%s%s ", WarnColor, nlog.WarnStr, ColorReset)
		cl.errorStr = fmt.Sprintf("%s%s%s ", ErrorColor, nlog.ErrorStr, ColorReset)
		cl.fatalStr = fmt.Sprintf("%s%s%s ", FatalColor, nlog.FatalStr, ColorReset)
	} else {
		cl.strW = strW
		cl.intW = intW
		cl.debugStr = fmt.Sprintf("%s ", nlog.DebugStr)
		cl.infoStr = fmt.Sprintf("%s ", nlog.InfoStr)
		cl.warnStr = fmt.Sprintf("%s ", nlog.WarnStr)
		cl.errorStr = fmt.Sprintf("%s ", nlog.ErrorStr)
		cl.fatalStr = fmt.Sprintf("%s ", nlog.FatalStr)
	}
}

// levelStr gets formatted level string
func (cl *Formatter) levelStr(lvl nlog.Level) string {
	switch lvl {
	case nlog.DEBUG:
		return cl.debugStr
	case nlog.WARNING:
		return cl.warnStr
	case nlog.ERROR:
		return cl.errorStr
	case nlog.FATAL:
		return cl.fatalStr
	}
	return cl.infoStr
}

func (cl *Formatter) hookSet(buff nlog.Buffer, lvl nlog.Level) nlog.HookFieldFn {
	return func(key string, value interface{}) {
		if cl.cfg.Level < lvl {
			return
		}
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

	cl.cfg.GetTime(buff, false)
	if buff.Len() > 0 {
		buff.AppendByte(' ')
	}
	// level
	if !cl.cfg.NoPrintLevel {
		buff.AppendString(cl.levelStr(lvl), false)
	}
	// logger Name
	if lgName != "" {
		buff.AppendString(fmt.Sprintf("%s[%s]%s", NameColor, lgName, ColorReset), false)
		buff.AppendByte(' ')
	}
	// Message
	if layout != "" {
		fmt.Fprintf(buff, layout, args...)
	} else {
		fmt.Fprint(buff, args...)
	}

	// Let assigning new fields from hooks even if fields is nil
	var hookBuff nlog.Buffer
	var buffSet nlog.HookBufferSet
	if len(cl.cfg.Hooks) > 0 {
		if fields == nil {
			hookBuff = pool.GetBuffer()
			defer pool.PutBuffer(hookBuff)
		} else {
			hookBuff = fields
		}
		buffSet = nlog.HookBufferSet{
			Buffer: hookBuff,
			With:   cl.hookSet(hookBuff, lvl),
		}
	}
	// Run hooks
	for i := range cl.cfg.Hooks {
		cl.cfg.Hooks[i].Run(lvl, buffSet, fmt.Sprintf(layout, args...))
	}
	// Fields filled only from hooks
	if hookBuff != nil {
		hookBuff.WriteTo(buff)
	}

	// Fields
	if fields != nil {
		buff.AppendByte(' ')
		fields.WriteTo(buff)
	}
	// File location
	if cl.cfg.FileLoc && LocColor != NoColor {
		buff.AppendByte(' ')
		if cl.cfg.Colored {
			buff.AppendString(LocColor, false)
		}
		formatter.GetFileLoc(cl.cfg.FileLocStrip, buff, callDepth, false)
		if cl.cfg.Colored {
			buff.AppendString(ColorReset, false)
		}
	}
	buff.AppendByte('\n')
	cl.cfg.Writer.WriteIfLevel(lvl, buff.Bytes())
}

// AppendFieldKey appends key value to log
func (cl *Formatter) AppendFieldKey(buf nlog.Buffer, key string) {
	buf.AppendByte(' ')
	if cl.cfg.Colored && KeyColor != NoColor {
		buf.AppendString(KeyColor, false)

	}
	buf.AppendString(key, false)
	if cl.cfg.Colored && KeyColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
	buf.AppendByte('=')
}

// AppendFieldValue appends value to log
func (cl *Formatter) AppendFieldValue(buf nlog.Buffer, val interface{}) {
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendAny(val, true, cl.cfg.MarshallFn)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// AppendKV formats and appends key,value to buffer
func (cl *Formatter) AppendKV(buf nlog.Buffer, lvl nlog.Level, key string, val interface{}) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendAny(val, false, cl.cfg.MarshallFn)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Str appends str value to buff with a format
func (cl *Formatter) Str(buf nlog.Buffer, lvl nlog.Level, key, val string) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendString(val, false)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Strs adds a slice of string value with a key to buff
func (cl *Formatter) Strs(buf nlog.Buffer, lvl nlog.Level, key string, val []string) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendStrings(val, false)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Int adds a new int key value to buff
func (cl *Formatter) Int(buf nlog.Buffer, lvl nlog.Level, key string, val int) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInt(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Ints adds a slice of int value with a key to buff
func (cl *Formatter) Ints(buf nlog.Buffer, lvl nlog.Level, key string, val []int) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInts(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Ints8 adds a slice of int8 value with a key to buff
func (cl *Formatter) Ints8(buf nlog.Buffer, lvl nlog.Level, key string, val []int8) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInts8(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Ints16 adds a slice of int16 value with a key to buff
func (cl *Formatter) Ints16(buf nlog.Buffer, lvl nlog.Level, key string, val []int16) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInts16(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Ints32 adds a slice of int32 value with a key to buff
func (cl *Formatter) Ints32(buf nlog.Buffer, lvl nlog.Level, key string, val []int32) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInts32(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Int64 adds a new int64 key value to buff
func (cl *Formatter) Int64(buf nlog.Buffer, lvl nlog.Level, key string, val int64) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInt64(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Int64s adds a slice of int64 value with a key to buff
func (cl *Formatter) Int64s(buf nlog.Buffer, lvl nlog.Level, key string, val []int64) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendInts64(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// UInt adds a new uint key value to buff
func (cl *Formatter) UInt(buf nlog.Buffer, lvl nlog.Level, key string, val uint) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendUInt64(uint64(val))
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// UInts adds a slice of uint value with a key to buff
func (cl *Formatter) UInts(buf nlog.Buffer, lvl nlog.Level, key string, val []uint) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendUInts(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// UInts16 adds a slice of uint16 value with a key to buff
func (cl *Formatter) UInts16(buf nlog.Buffer, lvl nlog.Level, key string, val []uint16) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendUInts16(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// UInts32 adds a slice of uint32 value with a key to buff
func (cl *Formatter) UInts32(buf nlog.Buffer, lvl nlog.Level, key string, val []uint32) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendUInts32(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Uint64 adds a new uint64 key value to buff
func (cl *Formatter) UInt64(buf nlog.Buffer, lvl nlog.Level, key string, val uint64) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendUInt64(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// UInts64 adds a slice of uint64 value with a key to buff
func (cl *Formatter) UInts64(buf nlog.Buffer, lvl nlog.Level, key string, val []uint64) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendUInts64(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Float32 adds a new float32 key value to buff
func (cl *Formatter) Float32(buf nlog.Buffer, lvl nlog.Level, key string, val float32) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendFloat32(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Floats32 adds a slice of float32 value with a key to buff
func (cl *Formatter) Floats32(buf nlog.Buffer, lvl nlog.Level, key string, val []float32) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendFloats32(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Float64 adds a new float64 key value to buff
func (cl *Formatter) Float64(buf nlog.Buffer, lvl nlog.Level, key string, val float64) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendFloat64(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Floats64 adds a slice of float32 value with a key to buff
func (cl *Formatter) Floats64(buf nlog.Buffer, lvl nlog.Level, key string, val []float64) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendFloats64(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Bool adds a new bool key value to buff
func (cl *Formatter) Bool(buf nlog.Buffer, lvl nlog.Level, key string, val bool) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendBool(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Bools adds a slice of bool value with a key to buff
func (cl *Formatter) Bools(buf nlog.Buffer, lvl nlog.Level, key string, val []bool) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendBools(val)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Error adds a new error key value to buff
func (cl *Formatter) Error(buf nlog.Buffer, lvl nlog.Level, key string, val error) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendError(val, false)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}

// Bools adds a slice of error value with a key to buff
func (cl *Formatter) Errors(buf nlog.Buffer, lvl nlog.Level, key string, val []error) {
	if cl.cfg.Level < lvl {
		return
	}
	cl.AppendFieldKey(buf, key)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ValColor, false)
	}
	buf.AppendErrors(val, false)
	if cl.cfg.Colored && ValColor != NoColor {
		buf.AppendString(ColorReset, false)
	}
}
