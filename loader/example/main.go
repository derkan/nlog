package main

import (
	"fmt"
	"os"

	"github.com/derkan/nlog/loader"
	"github.com/derkan/nlog/log"
)

func main() {
	filename := os.Args[1]
	l, err := loader.FromFile(filename, "log")
	if err != nil {
		fmt.Printf("Failed to load file '%s', err :%v\n", filename, err)
	}
	log.InitFromLoader(l, "nlogapp")
	log.Infof("test: %d", 123)
	log.Debug().Str("a", "string").Msg("my debug")
	log.Warn().Int("i", 4).Msg("my warning")
	log.Info().Bool("b", true).Strs("strslice", []string{"a"}).Msg("my bool")
	log.Error().Err(nil).With("k", map[string]int{"x": 1}).Msg("my err with map")
	log.Flush()
}
