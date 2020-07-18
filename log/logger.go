package log

import (
	"fmt"
	"os"

	"github.com/derkan/nlog"
	"github.com/derkan/nlog/formatter/console"
	"github.com/derkan/nlog/formatter/json"
	"github.com/derkan/nlog/loader"
	"github.com/derkan/nlog/pool"
)

// Init initalizes default logger
func Init(opts ...option) {
	Logger = New(opts...)
}

// InitFromLoader initalizes default logger
func InitFromLoader(cfg *loader.Loader, appName string) {
	Logger = NewFromConfig(cfg, appName)
}

// New returns a new instance of standard logger
func New(opts ...option) *Instance {
	ins := &Instance{cfg: &config{MinLevel: nlog.DEBUG}}
	for _, opt := range opts {
		opt(ins.cfg)
	}
	ins.SetDefaults()
	return ins
}

// NewFromConfig builds a logger from given loader config
// appName is used in syslog
func NewFromConfig(cfg *loader.Loader, appName string) (ins *Instance) {
	ins = &Instance{
		cfg: &config{
			MinLevel: cfg.Level,
			Prefix:   cfg.Prefix,
		},
	}
	for _, f := range cfg.Formatters {
		switch f.Type {
		case "console":
			ins.cfg.formatters = append(ins.cfg.formatters, console.NewFromConfig(f, appName))
		case "json":
			ins.cfg.formatters = append(ins.cfg.formatters, json.NewFromConfig(f, appName))
		default:
			fmt.Fprintf(os.Stderr, "invalid formater type %s\n", f.Type)
			continue
		}
	}
	ins.SetDefaults()
	return ins
}

func (ins *Instance) SetDefaults() {
	if len(ins.cfg.formatters) == 0 {
		ins.cfg.formatters = append(ins.cfg.formatters, console.NewFormatter(
			console.WithDate(),
			console.WithTime(),
		))
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Init()
	}
	ins.itemPool = pool.NewItemPool(4, nlog.DEBUG, ins.cfg.formatters...)
}

// Flush flushes to disk and closes writers
func (ins *Instance) Flush() {
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Flush()
	}
}

// Fatal returns FATAL level logger item
func (ins *Instance) Fatal() nlog.LoggerItem {
	if ins.cfg.MinLevel < nlog.FATAL {
		return pool.NullItem
	}
	return ins.itemPool.Get(nlog.FATAL, ins.cfg.Prefix, 0)
}

// Fatalf logs FATAL level log with given format-params and exits
func (ins *Instance) Fatalf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < nlog.FATAL {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), nlog.FATAL, nil, ins.cfg.Prefix, format, args...)
	}
}

// Error returns ERROR level logger item
func (ins *Instance) Error() nlog.LoggerItem {
	if ins.cfg.MinLevel < nlog.ERROR {
		return pool.NullItem
	}
	return ins.itemPool.Get(nlog.ERROR, ins.cfg.Prefix, 0)
}

// Errorf logs ERROR level log with given format-params
func (ins *Instance) Errorf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < nlog.ERROR {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), nlog.ERROR, nil, ins.cfg.Prefix, format, args...)
	}
}

// Warn returns WARNING level logger item
func (ins *Instance) Warn() nlog.LoggerItem {
	if ins.cfg.MinLevel < nlog.WARNING {
		return pool.NullItem
	}
	return ins.itemPool.Get(nlog.WARNING, ins.cfg.Prefix, 0)
}

// Warnf logs WARN level log with given format-params
func (ins *Instance) Warnf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < nlog.WARNING {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), nlog.WARNING, nil, ins.cfg.Prefix, format, args...)
	}
}

// Info logs info message if logging level is satisfied
func (ins *Instance) Info() nlog.LoggerItem {
	if ins.cfg.MinLevel < nlog.INFO {
		return pool.NullItem
	}
	return ins.itemPool.Get(nlog.INFO, ins.cfg.Prefix, 0)
}

// Infof logs info message with format if logging level is satisfied
func (ins *Instance) Infof(format string, args ...interface{}) {

	if ins.cfg.MinLevel < nlog.INFO {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), nlog.INFO, nil, ins.cfg.Prefix, format, args...)
	}
}

// Debug returns DEBUG level logger item
func (ins *Instance) Debug() nlog.LoggerItem {
	if ins.cfg.MinLevel < nlog.DEBUG {
		return pool.NullItem
	}
	return ins.itemPool.Get(nlog.DEBUG, ins.cfg.Prefix, 0)
}

// Debugf prints DEBUG level message with given format-params
func (ins *Instance) Debugf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < nlog.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, nlog.DEBUG, nil, ins.cfg.Prefix, format, args...)
	}
}

// Print prints log at INFO level with given params
func (ins *Instance) Print(args ...interface{}) {
	if ins.cfg.MinLevel < nlog.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, nlog.DEBUG, nil, ins.cfg.Prefix, "", args...)
	}
}

// Printf prints log at INFO level with given format-params
func (ins *Instance) Printf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < nlog.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, nlog.DEBUG, nil, ins.cfg.Prefix, format, args...)
	}
}

// Print prints log at INFO level with given params
func (ins *Instance) Println(args ...interface{}) {
	if ins.cfg.MinLevel < nlog.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, nlog.DEBUG, nil, ins.cfg.Prefix, "", args...)
	}
}

// Sub returns a sub logger with given prefix
func (ins *Instance) Sub(prefix string) interface{} {
	return &Instance{
		itemPool: ins.itemPool,
		cfg: &config{
			Prefix:     prefix,
			MinLevel:   ins.cfg.MinLevel,
			formatters: ins.cfg.formatters,
			SubDepth:   1,
		},
	}
}

// GetLevel returns logger level
func (ins *Instance) GetLevel() int {
	return int(ins.cfg.MinLevel)
}
