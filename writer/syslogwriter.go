// +build !windows

package writer

import (
	"fmt"
	"io"
	"log/syslog"
	"os"

	"github.com/derkan/nlog/common"
)

// SyslogWriter is an interface matching a syslog.Writer struct.
type SyslogWriter interface {
	io.Writer
	Debug(m string) error
	Info(m string) error
	Warning(m string) error
	Err(m string) error
	Emerg(m string) error
	Crit(m string) error
}

type syslogWriter struct {
	w SyslogWriter
	l common.Level
}

var logMap = map[common.Level]syslog.Priority{
	common.FATAL:   syslog.LOG_CRIT,
	common.ERROR:   syslog.LOG_ERR,
	common.WARNING: syslog.LOG_WARNING,
	common.INFO:    syslog.LOG_INFO,
	common.DEBUG:   syslog.LOG_DEBUG,
}

// SyslogLevelWriter wraps a SyslogWriter and call the right syslog level
// method matching the level.
func SysLogWrapper(w SyslogWriter, l common.Level) LeveledWriter {
	return syslogWriter{w, l}
}

func (sw syslogWriter) Write(p []byte) (n int, err error) {
	return sw.w.Write(p)
}

// Write implements io.WriteCloser.
func (l syslogWriter) Close() (err error) {
	return
}

// GetLevel returns log level of current writer
func (l syslogWriter) GetLevel() common.Level {
	return l.l
}

// WriteLevel implements LevelWriter interface.
func (sw syslogWriter) WriteIfLevel(lvl common.Level, p []byte) (n int, err error) {
	if lvl > sw.l {
		if p == nil {
			return 0, nil
		}
		return len(p), nil
	}
	switch lvl {
	case common.DEBUG:
		err = sw.w.Debug(string(p))
	case common.INFO:
		err = sw.w.Info(string(p))
	case common.WARNING:
		err = sw.w.Warning(string(p))
	case common.ERROR:
		err = sw.w.Err(string(p))
	case common.FATAL:
		err = sw.w.Crit(string(p))
	default:
		err = sw.w.Warning(string(p))
	}
	n = len(p)
	return
}

// NewSysLogWriter initializes a syslog writer wrapped in LeveledLogger
func NewSysLogWriter(name string, l ...common.Level) LeveledWriter {
	lvl := common.DEBUG
	if len(l) > 0 {
		lvl = l[0]
	}
	w, err := syslog.New(logMap[lvl], name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Check if rsyslog service is running, err: %v\n", err)
		return NewDummyLeveledWriter()
	}
	return SysLogWrapper(w, lvl)
}
