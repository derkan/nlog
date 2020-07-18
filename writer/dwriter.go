package writer

import "github.com/derkan/nlog/common"

// DummyLeveledWriter is a wrapper around an actual writer
// Do not do anything
type DummyLeveledWriter struct {
}

// Close closes writer
func (l *DummyLeveledWriter) Close() error {
	return nil
}

// Write Writes data
func (l *DummyLeveledWriter) Write(p []byte) (n int, err error) {
	return
}

// WriteIfLevel writes mesage if level is satisfied
func (l *DummyLeveledWriter) WriteIfLevel(lvl common.Level, p []byte) (n int, err error) {
	return
}

// GetLevel returns log level
func (l *DummyLeveledWriter) GetLevel() common.Level {
	return common.FATAL
}

// NewDummyLeveledWriter returns a new instance of DummyLeveledWriter
// Do not write anything
func NewDummyLeveledWriter() *DummyLeveledWriter {
	return new(DummyLeveledWriter)
}
