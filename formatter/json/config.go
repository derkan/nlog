package json

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/derkan/nlog/common"
	"github.com/derkan/nlog/writer"
)

// config is for logger settings
type config struct {
	// Level is logging level
	Level common.Level `json:"level" yaml:"level"`
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
	// Hooks hold hook structs to be called during logging
	Hooks []common.Hook
	// MarshallFn func is used to serialize interfaces
	MarshallFn common.MarshallFn
}

// option is function type used for setting console logging config attributes
type option func(*config)

// WithMarshallFn adds marshalling with specified func
// default is json.Marshall
func WithMarshallFn(fn common.MarshallFn) option {
	return func(c *config) {
		c.MarshallFn = fn
	}
}

// WithHook adds hooks to be called during logging
func WithHook(h common.Hook) option {
	return func(c *config) {
		c.Hooks = append(c.Hooks, h)
	}
}

// WithWriter sets the Writer for console
// Can be called mutliple times to add new writers
func WithWriter(w io.WriteCloser, lvl ...common.Level) option {
	return func(c *config) {
		var wl *writer.Writer
		if len(lvl) == 0 {
			wl = writer.NewWriter(w, common.INFO)
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
func WithParallelWriter(w io.WriteCloser, chanSize int, lvl ...common.Level) option {
	return func(c *config) {
		var wl *writer.ParallelWriter
		if len(lvl) == 0 {
			wl = writer.NewParellelWriter(w, common.INFO, chanSize)
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
func WithLevel(level common.Level) option {
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
func (c *config) GetTime(buf common.Buffer, quota bool) {
	if !c.Date && !c.Time && !c.UnixTime {
		return
	}
	fmt.Fprintf(buf, "%q:", TimeKey)
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
			buf.AppendByte(' ')
		}

		if c.Time {
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
			if c.TimeResolution < time.Millisecond {
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
