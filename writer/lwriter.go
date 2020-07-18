package writer

import (
	"io"

	"github.com/derkan/nlog"
)

// LeveledWriter is a wrapper around an actual writer. Not threadsafe
type LeveledWriter interface {
	io.WriteCloser
	WriteIfLevel(lvl nlog.Level, p []byte) (n int, err error)
	GetLevel() nlog.Level
}

// Writer is writer with log level filtering
// It wraps real writer and calls it for logs where logging level is satisfied
// Not concurrent safe
type Writer struct {
	w io.WriteCloser
	l nlog.Level
}

// Write implements io.Writer.
func (l *Writer) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}

// Write implements io.WriteCloser.
func (l *Writer) Close() (err error) {
	return l.w.Close()
}

// GetLevel returns log level of current writer
func (l *Writer) GetLevel() nlog.Level {
	return l.l
}

// WriteIfLevel calls write if current le₺vel is satisfied
func (l *Writer) WriteIfLevel(lvl nlog.Level, p []byte) (n int, err error) {
	if lvl > l.l {
		if p == nil {
			return 0, nil
		}
		return len(p), nil
	}
	if lw, ok := l.w.(LeveledWriter); ok {
		return lw.WriteIfLevel(lvl, p)
	}
	return l.Write(p)
}

// NewWriter returns a new instance of threadsafe writer which will
// write when level is satisfied
func NewWriter(w io.WriteCloser, l nlog.Level) *Writer {
	return &Writer{
		w: w,
		l: l,
	}
}
