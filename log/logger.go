package log

import (
	"fmt"
	"os"

	"github.com/derkan/nlog/common"
	"github.com/derkan/nlog/formatter/console"
	"github.com/derkan/nlog/formatter/json"
	"github.com/derkan/nlog/loader"
	"github.com/derkan/nlog/pool"
)

// Init initalizes default logger
func Init(opts ...option) {
	StdLogger = New(opts...)
}

// InitFromLoader initalizes default logger
func InitFromLoader(cfg *loader.Loader, appName string) {
	StdLogger = NewFromConfig(cfg, appName)
}

// New returns a new instance of standard logger
func New(opts ...option) *Instance {
	ins := &Instance{cfg: &config{MinLevel: common.DEBUG}}
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
			MinLevel: cfg.MinLevel,
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
	ins.itemPool = pool.NewItemPool(4, common.DEBUG, ins.cfg.formatters...)
}

// Flush flushes to disk and closes writers
func (ins *Instance) Flush() {
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Flush()
	}
}

// Fatal returns FATAL level logger item
func (ins *Instance) Fatal() common.LoggerItem {
	if ins.cfg.MinLevel < common.FATAL {
		return pool.NullItem
	}
	return ins.itemPool.Get(common.FATAL, ins.cfg.Prefix, 0)
}

// Fatalf logs FATAL level log with given format-params and exits
func (ins *Instance) Fatalf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < common.FATAL {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), common.FATAL, nil, ins.cfg.Prefix, format, args...)
	}
}

// Error returns ERROR level logger item
func (ins *Instance) Error() common.LoggerItem {
	if ins.cfg.MinLevel < common.ERROR {
		return pool.NullItem
	}
	return ins.itemPool.Get(common.ERROR, ins.cfg.Prefix, 0)
}

// Errorf logs ERROR level log with given format-params
func (ins *Instance) Errorf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < common.ERROR {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), common.ERROR, nil, ins.cfg.Prefix, format, args...)
	}
}

// Warn returns WARNING level logger item
func (ins *Instance) Warn() common.LoggerItem {
	if ins.cfg.MinLevel < common.WARNING {
		return pool.NullItem
	}
	return ins.itemPool.Get(common.WARNING, ins.cfg.Prefix, 0)
}

// Warnf logs WARN level log with given format-params
func (ins *Instance) Warnf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < common.WARNING {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), common.WARNING, nil, ins.cfg.Prefix, format, args...)
	}
}

// Info logs info message if logging level is satisfied
func (ins *Instance) Info() common.LoggerItem {
	if ins.cfg.MinLevel < common.INFO {
		return pool.NullItem
	}
	return ins.itemPool.Get(common.INFO, ins.cfg.Prefix, 0)
}

// Infof logs info message with format if logging level is satisfied
func (ins *Instance) Infof(format string, args ...interface{}) {

	if ins.cfg.MinLevel < common.INFO {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(ins.cfg.formatters[i].GetCallDepth(ins.cfg.SubDepth), common.INFO, nil, ins.cfg.Prefix, format, args...)
	}
}

// Debug returns DEBUG level logger item
func (ins *Instance) Debug() common.LoggerItem {
	if ins.cfg.MinLevel < common.DEBUG {
		return pool.NullItem
	}
	return ins.itemPool.Get(common.DEBUG, ins.cfg.Prefix, 0)
}

// Debugf prints DEBUG level message with given format-params
func (ins *Instance) Debugf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < common.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, common.DEBUG, nil, ins.cfg.Prefix, format, args...)
	}
}

// Print prints log at INFO level with given params
func (ins *Instance) Print(args ...interface{}) {
	if ins.cfg.MinLevel < common.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, common.DEBUG, nil, ins.cfg.Prefix, "", args...)
	}
}

// Printf prints log at INFO level with given format-params
func (ins *Instance) Printf(format string, args ...interface{}) {
	if ins.cfg.MinLevel < common.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, common.DEBUG, nil, ins.cfg.Prefix, format, args...)
	}
}

// Print prints log at INFO level with given params
func (ins *Instance) Println(args ...interface{}) {
	if ins.cfg.MinLevel < common.DEBUG {
		return
	}
	for i := range ins.cfg.formatters {
		ins.cfg.formatters[i].Logf(0, common.DEBUG, nil, ins.cfg.Prefix, "", args...)
	}
}

// Sub returns a sub logger with given prefix and optionally min logging level
func (ins *Instance) Sub(prefix string, minLevel ...common.Level) Logger {
	lvl := ins.cfg.MinLevel
	if len(minLevel) > 0 {
		lvl = minLevel[0]
	}
	return &Instance{
		itemPool: ins.itemPool,
		cfg: &config{
			Prefix:     prefix,
			MinLevel:   lvl,
			formatters: ins.cfg.formatters,
			SubDepth:   1,
		},
	}
}
