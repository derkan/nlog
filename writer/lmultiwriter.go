package writer

import (
	"fmt"
	"io"
	"sync"

	"github.com/derkan/nlog/common"
)

// LeveledWriter is a wrapper around an actual writer. Not threadsafe
type LeveledMultiWriter interface {
	io.WriteCloser
	WriteIfLevel(lvl common.Level, p []byte) (n int, err error)
	Remove(writers ...LeveledWriter)
	Append(writers ...LeveledWriter)
}

// MultiWriter is multi writer with log level filtering
// It wraps real writer and calls it for logs where logging level is satisfied
// Not concurrent safe
type MultiWriter struct {
	mu      sync.Mutex
	writers []*Writer
}

func (t *MultiWriter) Remove(writers ...LeveledWriter) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := len(t.writers) - 1; i > 0; i-- {
		for _, v := range writers {
			if t.writers[i] == v {
				t.writers = append(t.writers[:i], t.writers[i+1:]...)
				break
			}
		}
	}
}
func (t *MultiWriter) Append(writers ...LeveledWriter) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, w := range writers {
		if wx, ok := w.(*Writer); ok {
			t.writers = append(t.writers, wx)
		}
	}
}

// Write implements io.Writer.
func (t *MultiWriter) Write(p []byte) (n int, err error) {
	var werr error
	for i, w := range t.writers {
		if n, werr = w.Write(p); err != nil {
			err = fmt.Errorf("%v, worker: %d, err: %v", i, err, werr)
		}
		if n != len(p) {
			err = fmt.Errorf("%v, worker: %d, err: %v", i, err, io.ErrShortWrite)
			return
		}
	}
	return len(p), nil
}

// Close implements io.WriteCloser.
func (t *MultiWriter) Close() (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i := range t.writers {
		t.writers[i].Close()
	}
	return
}

// WriteIfLevel calls write if current le₺vel is satisfied
func (t *MultiWriter) WriteIfLevel(lvl common.Level, p []byte) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, w := range t.writers {
		if lvl > w.GetLevel() {
			continue
		}
		if _, werr := w.WriteIfLevel(lvl, p); err != nil {
			err = fmt.Errorf("%v, worker: %d, err: %v", i, err, werr)
		}
	}
	return len(p), err
}

// NewMultiWriter returns a new instance of threadsafe writer which will
// write when level is satisfied
func NewMultiWriter(writers ...*Writer) *MultiWriter {
	w := make([]*Writer, len(writers))
	copy(w, writers)
	return &MultiWriter{writers: w}
}
