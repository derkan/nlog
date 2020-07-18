package writer

import (
	"io"
	"sync"
	"time"

	"github.com/derkan/nlog/common"
	"github.com/derkan/nlog/pool"
)

// ParallelWriter is writer with log level filtering
// It wraps real writer and calls it for logs where logging level is satisfied
// Not concurrent safe
type ParallelWriter struct {
	w   io.WriteCloser
	l   common.Level
	ch  chan common.Buffer
	end chan bool
}

// Write implements io.Writer.
func (l *ParallelWriter) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}

// Write implements io.WriteCloser.
func (l *ParallelWriter) Close() (err error) {
	return l.w.Close()
}

// GetLevel returns log level of current writer
func (l *ParallelWriter) GetLevel() common.Level {
	return l.l
}

// WriteIfLevel calls write if current le₺vel is satisfied
func (l *ParallelWriter) WriteIfLevel(lvl common.Level, p []byte) (n int, err error) {
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

// Start starts worker for writer
func (l *ParallelWriter) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for {
			select {
			case p := <-l.ch: // worker has received job
				l.Write(p.Bytes())
				pool.PutBuffer(p)
			case <-l.end:
				defer wg.Done()
				// wait until all data in channel is written or timeout happens before exiting:
				for len(l.ch) > 0 {
					select {
					case p := <-l.ch:
						l.Write(p.Bytes())
						pool.PutBuffer(p)
					case <-time.After(2):
						break
					}
				}
				l.Close()
				return
			}
		}
	}()
}

// Stop stops worker for writer
func (l *ParallelWriter) Stop() {
	l.end <- true
}

// NewParellelWriter returns a new instance of threadsafe writer which will
// write when level is satisfied
func NewParellelWriter(w io.WriteCloser, l common.Level, chanSize int) *ParallelWriter {
	return &ParallelWriter{
		w:   w,
		l:   l,
		ch:  make(chan common.Buffer, chanSize),
		end: make(chan bool),
	}
}
