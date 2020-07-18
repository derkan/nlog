package main

import (
	"sync"
	"time"

	"github.com/derkan/nlog/common"
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
		log.WithMinLevel(common.DEBUG),
		log.WithFormatter(
			console.NewFormatter(
				console.WithColor(),
				console.WithDate(),
				console.WithTime(time.Millisecond),
				console.WithFileLoc(),
				console.WithStripPath("/data/go/src"),
				console.WithLevel(common.DEBUG),
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
				console.WithParallelWriter(writer.NewSysLogWriter("MYAPP"), 100, common.DEBUG),
			),
		))
	log.Infof("starting")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		time.Sleep(0)
		for i := 0; i < 1000; i++ {
			log.Infof("test: %d", 123)
		}
		wg.Done()
	}()

	subLog := log.Sub("wapp")
	go func() {
		for i := 0; i < 1000; i++ {
			subLog.Warnf("test: %d", 123)
		}
		wg.Done()
	}()
	wg.Wait()
	subLog.Infof("done")
	log.Flush()
}
