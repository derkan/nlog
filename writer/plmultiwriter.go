package writer

import (
	"fmt"
	"io"
	"sync"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/pool"
)

// ParallelMultiWriter is multi writer with log level filtering
// It wraps real writer and calls it for logs where logging level is satisfied
// Not concurrent safe
type ParallelMultiWriter struct {
	writers []*ParallelWriter
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

func (t *ParallelMultiWriter) Remove(writers ...LeveledWriter) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := len(t.writers) - 1; i > 0; i-- {
		for _, v := range writers {
			if t.writers[i] == v {
				t.writers[i].Stop()
				t.writers = append(t.writers[:i], t.writers[i+1:]...)
				break
			}
		}
	}
}
func (t *ParallelMultiWriter) Append(writers ...LeveledWriter) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, w := range writers {
		if wx, ok := w.(*ParallelWriter); ok {
			t.writers = append(t.writers, wx)
			wx.Start(&t.wg)
		}
	}
}

// Write implements io.Writer.
func (t *ParallelMultiWriter) Write(p []byte) (n int, err error) {
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
func (t *ParallelMultiWriter) Close() (err error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, w := range t.writers {
		w.Stop() // stop worker
	}
	t.wg.Wait()
	return
}

// WriteIfLevel calls write if current le₺vel is satisfied
func (t *ParallelMultiWriter) WriteIfLevel(lvl nlog.Level, p []byte) (n int, err error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, w := range t.writers {
		if lvl > w.GetLevel() {
			continue
		}
		w.ch <- pool.GetBuffer().AppendBytes(p)

	}
	return len(p), err
}

// NewParallelMultiWriter returns a new instance of threadsafe writer which will
// write when level is satisfied
func NewParallelMultiWriter(writers ...*ParallelWriter) *ParallelMultiWriter {
	w := make([]*ParallelWriter, len(writers))
	copy(w, writers)
	res := &ParallelMultiWriter{writers: w}
	for i := range w {
		w[i].Start(&res.wg)
	}
	return res
}
