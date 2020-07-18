package log

import (
	"github.com/derkan/nlog/common"
	"github.com/derkan/nlog/pool"
)

// config is for logger settings
type config struct {
	// Prefix is logger prefix, this prefix is printed before message on each line
	Prefix string
	// MinLevel is minimum log level permitted for writers to write
	MinLevel   common.Level
	formatters []common.Formatter
	SubDepth   int
}

// option is function type used for setting Config
type option func(*config)

// Instance is logger instance
type Instance struct {
	cfg      *config
	itemPool pool.ItemPool
}

// WithPrefix sets Prefix for logger
// Prefix is shown after log level
func WithPrefix(Prefix string) option {
	return func(c *config) {
		if Prefix != "" {
			c.Prefix = Prefix
		}
	}
}

// WithMinLevel sets minimal level of logger
func WithMinLevel(level common.Level) option {
	return func(c *config) {
		c.MinLevel = level
	}
}

// WithFormatter adds formatter to  logger
func WithFormatter(formatter common.Formatter) option {
	return func(c *config) {
		c.formatters = append(c.formatters, formatter)
	}
}
