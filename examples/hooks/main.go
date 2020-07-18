package main

import (
	"time"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/formatter/json"
	"github.com/derkan/nlog/log"
)

// jsoniter "github.com/json-iterator/go"
// var json = jsoniter.ConfigCompatibleWithStandardLibrary
func main() {
	log.Init(
		//log.WithMinLevel(nlog.INFO),
		log.WithFormatter(
			json.NewFormatter(
				json.WithLevel(nlog.INFO),
				json.WithDate(),
				json.WithTime(time.Millisecond),
				json.WithHook(nlog.HookFunc(func(level nlog.Level, hSet nlog.HookBufferSet, msg string) {
					hSet.With("filled_from_hook", true)
				})),
			),
		))
	log.Infof("test: %d", 123)
	log.Debug().Str("a", "string").Msg("my debug")
	log.Warn().Int("i", 4).Msg("my warning")
	log.Info().Bool("b", true).Strs("strslice", []string{"a"}).Msg("my bool")
	log.Error().Err(nil).With("k", map[string]int{"x": 1}).Msg("my err with map")
	log.Flush()
}
