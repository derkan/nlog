package console

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/writer"
)

var (
	// ColorReset resets color code to default
	ColorReset string = "\033[0m"
	//NoColor is used to dismiss coloring
	NoColor string = ""
	// Black is ANSI color code for black
	Black string = "\033[1;30m"
	// Red is ANSI color code for red
	Red string = "\033[1;31m"
	// Green is ANSI color code for green
	Green string = "\033[1;32m"
	// Yellow is ANSI color code for yellow
	Yellow string = "\033[1;33m"
	// Blue is ANSI color code for blue
	Blue string = "\033[1;34m"
	// Magenta is ANSI color code for magenta
	Magenta string = "\033[1;35m"
	// Cyan is ANSI color code for cyan
	Cyan string = "\033[1;36m"
	// White is ANSI color code for white
	White string = "\033[1;37m"

	// BlackBold is ANSI color code for bold black
	BlackBold string = "\033[1;30;1m"
	// RedBold is ANSI color code for bold red
	RedBold string = "\033[1;31;1m"
	// GreenBold is ANSI color code for bold green
	GreenBold string = "\033[1;32;1m"
	// YellowBold is ANSI color code for bold yellow
	YellowBold string = "\033[1;33;1m"
	// BlueBold is ANSI color code for bold blue
	BlueBold string = "\033[1;34;1m"
	// MagentaBold is ANSI color code for bold magenta
	MagentaBold string = "\033[1;35;1m"
	// CyanBold is ANSI color code for bold cyan
	CyanBold string = "\033[1;36;1m"
	// WhiteBold is ANSI color code for bold white
	WhiteBold string = "\033[1;37;1m"
)

var (
	// DebugColor is color for debug keyword
	DebugColor string = Green
	// InfoColor  is color for info keyword
	InfoColor string = Cyan
	// NoticeColor  is color for notice keyword
	NoticeColor string = White
	// WarnColor  is color for warn keyword
	WarnColor string = Yellow
	// ErrorColor  is color for error keyword
	ErrorColor string = Red
	// FatalColor  is color for fatal keyword
	FatalColor string = Magenta

	// KeyColor  is color for keys in key=value
	KeyColor = Magenta
	// ValColor  is color for values in key=value
	ValColor = Cyan
	// LocColor  is file:line color
	LocColor = Black
	// NameColor  is logger prefix color
	NameColor = GreenBold
	// TimeColor  is logging time color
	TimeColor = NoColor
)

// LevelColors holds colors for each log level
var LevelColors = map[nlog.Level]string{
	nlog.FATAL:   FatalColor,
	nlog.ERROR:   ErrorColor,
	nlog.WARNING: WarnColor,
	nlog.INFO:    InfoColor,
	nlog.DEBUG:   DebugColor,
}

// config is for logger settings
type config struct {
	// Level is logging level
	Level nlog.Level `json:"level" yaml:"level"`
	// NoPrintLevel if set level string will not be written
	NoPrintLevel bool `json:"no_print_level" yaml:"print_level"`
	// Date sets whether to print date or not in layout 2006-01-02
	Date bool `json:"date" yaml:"date"`
	// Whether to print as string time, resolution can be set using TimeResolution
	Time bool `json:"time" yaml:"time"`
	// Prints time in UTC if defined
	TimeUTC bool `json:"time_utc" yaml:"time_utc"`
	// Whether to print in unix time (as int), resolution can be set using TimeResolution
	UnixTime bool `json:"unix_time" yaml:"unix_time"`
	// TimeResolution is used to getting time with specific resolution
	TimeResolution time.Duration `json:"time_resolution" yaml:"time_resolution"`
	// FileLoc is for printing caller file location or not
	FileLoc bool `json:"file_loc" yaml:"file_loc"`
	// FileLocStrip is for stripping given path prefix from file location
	FileLocStrip string `json:"file_loc_strip" yaml:"file_loc_strip"`
	// FileLocCallerDepth is for setting depth of caller to find file loc
	FileLocCallerDepth int `json:"file_loc_caller_depth" yaml:"file_loc_caller_depth"`
	// Writer is writer to write log into
	Writer writer.LeveledMultiWriter
	// MarshallFn func is used to serialize interfaces
	MarshallFn nlog.MarshallFn
	// Hooks hold hook structs to be called during logging
	Hooks []nlog.Hook
	// Colored determines whether to print in color
	Colored bool
}

// option is function type used for setting console logging config attributes
type option func(*config)

// WithColor adds color support
func WithColor() option {
	return func(c *config) {
		c.Colored = true
	}
}

// WithMarshallFn adds marshalling with specified func
// default is json.Marshall
func WithMarshallFn(fn nlog.MarshallFn) option {
	return func(c *config) {
		c.MarshallFn = fn
	}
}

// WithHook adds hooks to be called during logging
func WithHook(h nlog.Hook) option {
	return func(c *config) {
		c.Hooks = append(c.Hooks, h)
	}
}

// WithWriter sets the Writer for console
// Can be called mutliple times to add new writers
// Overrides if called after WithParallelWriter
func WithWriter(w io.WriteCloser, lvl ...nlog.Level) option {
	return func(c *config) {
		var wl *writer.Writer
		if len(lvl) == 0 {
			wl = writer.NewWriter(w, nlog.INFO)
		} else {
			wl = writer.NewWriter(w, lvl[0])
		}
		if c.Writer == nil {
			c.Writer = writer.NewMultiWriter(wl)
		} else {
			c.Writer.Append(wl)
		}
	}
}

// WithParallelWriter sets the Writer for console
// Can be called mutliple times to add new writers
// Overrides if called after WithWriter
func WithParallelWriter(w io.WriteCloser, chanSize int, lvl ...nlog.Level) option {
	return func(c *config) {
		var wl *writer.ParallelWriter
		if len(lvl) == 0 {
			wl = writer.NewParellelWriter(w, nlog.INFO, chanSize)
		} else {
			wl = writer.NewParellelWriter(w, lvl[0], chanSize)
		}
		if c.Writer == nil {
			c.Writer = writer.NewParallelMultiWriter(wl)
		} else {
			c.Writer.Append(wl)
		}
	}
}

// WithLevel sets level of logger
func WithLevel(level nlog.Level) option {
	return func(c *config) {
		c.Level = level
	}
}

// WithNoPrintLevel prints date
func WithNoPrintLevel() option {
	return func(c *config) {
		c.NoPrintLevel = true
	}
}

// WithTime adds time to log with optional resolution. Default resolution is time.Second
// Overrides if WithUnixTime if called after it
func WithTime(resolution ...time.Duration) option {
	return func(c *config) {
		c.Time = true
		if len(resolution) > 0 {
			c.TimeResolution = resolution[0]
		} else {
			c.TimeResolution = time.Second
		}
	}
}

// WithUnixTime prints time in unix time format
// Resolution may be set in WithTimeNano, WithTimeMicro, WithTimeMilli
// Overrides if WithTime if called after it
func WithUnixTime(resolution ...time.Duration) option {
	return func(c *config) {
		c.UnixTime = true
		if len(resolution) > 0 {
			c.TimeResolution = resolution[0]
		} else {
			c.TimeResolution = time.Second
		}

	}
}

// WithTimeUTC prints time in UTC
func WithTimeUTC() option {
	return func(c *config) {
		c.TimeUTC = true
	}
}

// WithDate prints date
func WithDate() option {
	return func(c *config) {
		c.Date = true
	}
}

// WithStripPath strips given Prefix from file path
func WithStripPath(Prefix string) option {
	addSlash := func(s string) string {
		if !strings.HasSuffix(s, string(os.PathSeparator)) {
			return fmt.Sprintf("%s%c", s, os.PathSeparator)
		}
		return s
	}
	return func(c *config) {
		c.FileLocStrip = filepath.Dir(addSlash(Prefix))
		if !strings.HasSuffix(c.FileLocStrip, string(os.PathSeparator)) {
			c.FileLocStrip = fmt.Sprintf("%s%c", c.FileLocStrip, os.PathSeparator)
		}
	}
}

// WithFileLoc prints file:line
func WithFileLoc(CallerDepth ...int) option {
	return func(c *config) {
		if len(CallerDepth) > 0 && CallerDepth[0] > 0 {
			c.FileLocCallerDepth = CallerDepth[0]
		}
		c.FileLoc = true
	}
}

// GetTime formats time
func (c *config) GetTime(buf nlog.Buffer, quota bool) {
	if !c.Date && !c.Time && !c.UnixTime {
		return
	}
	t := time.Now()
	if c.UnixTime {
		buf.AppendInt64(t.UnixNano() / int64(c.TimeResolution))
	} else {
		if c.TimeUTC {
			t = t.UTC()
		}
		if quota {
			buf.AppendByte('"')
		}
		if c.Date {
			year, month, day := t.Date()
			buf.Itoa(year, 4)
			buf.AppendByte('/')
			buf.Itoa(int(month), 2)
			buf.AppendByte('/')
			buf.Itoa(day, 2)
		}

		if c.Time {
			buf.AppendByte(' ')
			hour, min, sec := t.Clock()
			if c.TimeResolution <= time.Hour {
				buf.Itoa(hour, 2)
				buf.AppendByte(':')
			}
			if c.TimeResolution <= time.Minute {
				buf.Itoa(min, 2)
				buf.AppendByte(':')
			}
			if c.TimeResolution <= time.Second {
				buf.Itoa(sec, 2)
			}
			if c.TimeResolution < time.Second {
				buf.AppendByte('.')
			}
			if c.TimeResolution == time.Millisecond {
				buf.Itoa(t.Nanosecond()/1e3, 6)
			} else if c.TimeResolution == time.Microsecond {
				buf.Itoa(t.Nanosecond()/1e3, 3)
			} else if c.TimeResolution == time.Nanosecond {
				buf.Itoa(t.Nanosecond(), 9)
			}
		}
		if quota {
			buf.AppendByte('"')
		}
	}
}
