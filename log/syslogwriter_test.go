// +build !windows

package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/derkan/nlog/common"
	"github.com/derkan/nlog/formatter/json"
	"github.com/derkan/nlog/writer"
)

type syslogEvent struct {
	level string
	msg   string
}
type syslogTestWriter struct {
	events []syslogEvent
}

func (w *syslogTestWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func (w *syslogTestWriter) WriteIfLevel(p []byte, l common.Level) (int, error) {
	return 0, nil
}

func (w *syslogTestWriter) Trace(m string) error {
	w.events = append(w.events, syslogEvent{"Trace", m})
	return nil
}
func (w *syslogTestWriter) Debug(m string) error {
	fmt.Printf("--- debug:%s\n", m)
	w.events = append(w.events, syslogEvent{"Debug", m})
	return nil
}
func (w *syslogTestWriter) Info(m string) error {
	w.events = append(w.events, syslogEvent{"Info", m})
	return nil
}
func (w *syslogTestWriter) Warning(m string) error {
	w.events = append(w.events, syslogEvent{"Warning", m})
	return nil
}
func (w *syslogTestWriter) Err(m string) error {
	w.events = append(w.events, syslogEvent{"Err", m})
	return nil
}
func (w *syslogTestWriter) Emerg(m string) error {
	w.events = append(w.events, syslogEvent{"Emerg", m})
	return nil
}
func (w *syslogTestWriter) Crit(m string) error {
	w.events = append(w.events, syslogEvent{"Crit", m})
	return nil
}

func TestSyslogWriter(t *testing.T) {
	sw := &syslogTestWriter{}
	log := New(
		WithMinLevel(common.DEBUG),
		WithFormatter(
			json.NewFormatter(
				json.WithWriter(writer.SysLogWrapper(sw, common.DEBUG), common.DEBUG),
			),
		),
	)
	log.Debug().Msg("debug")
	log.Info().Msg("info")
	log.Warn().Msg("warn")
	log.Error().Msg("error")
	want := []syslogEvent{
		{"Debug", `{"level":"DBG","msg":"debug"}` + "\n"},
		{"Info", `{"level":"INF","msg":"info"}` + "\n"},
		{"Warning", `{"level":"WRN","msg":"warn"}` + "\n"},
		{"Err", `{"level":"ERR","msg":"error"}` + "\n"},
	}
	if got := sw.events; !reflect.DeepEqual(got, want) {
		t.Errorf("Invalid syslog message routing: want %v, got %v", want, got)
	}
}
