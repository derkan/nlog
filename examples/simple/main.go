package main

import (
	"time"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/formatter/console"
	"github.com/derkan/nlog/formatter/json"
	"github.com/derkan/nlog/log"
	"github.com/derkan/nlog/writer"
	fl "github.com/derkan/nlog/writer/filerotater"
)

// jsoniter "github.com/json-iterator/go"
// var json = jsoniter.ConfigCompatibleWithStandardLibrary
func main() {
	log.Init(
		log.WithPrefix("main"),
		log.WithMinLevel(nlog.DEBUG),
		log.WithFormatter(
			console.NewFormatter(
				console.WithColor(),
				console.WithDate(),
				console.WithTime(time.Millisecond),
				console.WithFileLoc(),
				console.WithStripPath("/data/go/src/swarmdb"),
				console.WithLevel(nlog.DEBUG),
			),
		), log.WithFormatter(
			json.NewFormatter(
				json.WithDate(),
				json.WithTime(),
				json.WithParallelWriter(
					fl.NewFileRotater(
						fl.WithFilename("/tmp/app.log"),
						fl.WithMaxBackups(3),
						fl.WithCompress(),
					), 100),
			)), log.WithFormatter(
			console.NewFormatter(
				console.WithParallelWriter(writer.NewSysLogWriter("MYAPP"), 100, nlog.DEBUG),
			),
		))
	log.Infof("test: %d", 123)
	log.Debug().Str("a", "string").Msg("my debug")
	log.Warn().Int("i", 4).Msg("my warning")
	log.Info().Bool("b", true).Strs("strslice", []string{"a"}).Msg("my bool")
	log.Error().Err(nil).With("k", map[string]int{"x": 1}).Msg("my err with map")

	subLog := log.Sub("wapp")
	subLog.Infof("sub infof: %d", 123)
	subLog.Warn().Int("i", 4).Msg("sub .Warn() warning")
	log.Infof("test: %d", 124)
	log.Flush()
}
