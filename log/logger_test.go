package log

import (
	"bytes"
	"testing"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/formatter/console"
)

type BuffCloser struct {
	*bytes.Buffer
}

func (b *BuffCloser) Close() error {
	// Noop
	return nil
}

/*
go test -v -bench=. -benchmem  -memprofile /tmp/profile_mem.out
go tool pprof -svg /tmp/profile_mem.out > /tmp/profile_mem.svg
firefox /tmp/profile_mem.svg
*/
// go test  -bench=Infof -benchmem
func BenchmarkInfof(b *testing.B) {
	output := &BuffCloser{&bytes.Buffer{}}
	log := New(
		WithMinLevel(nlog.DEBUG),
		WithFormatter(
			console.NewFormatter(
				console.WithColor(),
				console.WithDate(),
				console.WithTime(), //console.WithFileLoc(),
				console.WithWriter(output, nlog.DEBUG),
			),
		),
	)
	//BenchmarkInfof-4          512599              2077 ns/op             715 B/op          3 allocs/op

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		log.Infof("a string %s", "test")
	}
}

// go test  -bench=Str -benchmem
func BenchmarkStr(b *testing.B) {
	output := &BuffCloser{&bytes.Buffer{}}
	log := New(
		WithMinLevel(nlog.DEBUG),
		//log.WithFormatter(&log.JSONFormatter{}),
		WithFormatter(
			console.NewFormatter(
				console.WithColor(),
				console.WithDate(),
				console.WithTime(), //console.WithFileLoc(),
				console.WithWriter(output, nlog.DEBUG),
			),
		),
	)
	// BenchmarkStr-4            533245              1926 ns/op             744 B/op          3 allocs/op

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		log.Info().Str("a", "string").Msg("test")
	}
}

// go test  -bench=With -benchmem
func BenchmarkWith(b *testing.B) {
	output := &BuffCloser{&bytes.Buffer{}}
	log := New(
		WithMinLevel(nlog.DEBUG),
		WithFormatter(
			console.NewFormatter(
				console.WithColor(),
				console.WithDate(),
				console.WithTime(), //console.WithFileLoc(),
				console.WithWriter(output, nlog.DEBUG),
			),
		),
	)
	// BenchmarkStr-4            533245              1926 ns/op             744 B/op          3 allocs/op

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		log.Info().With("a", "string").Msg("test")
	}
}

/*
func BenchmarkZeroInfof(b *testing.B) {
	//	"github.com/rs/zerolog"
	//	zlog "github.com/rs/zerolog/log"

	output := &bytes.Buffer{}
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339Nano, NoColor: false}).With().Caller().Logger()

	// BenchmarkZeroInfof-4       48831             24805 ns/op            4639 B/op        102 allocs/op

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		zlog.Logger.Info().Str("a", "string").Msg("test")
	}
}
*/
