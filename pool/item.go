package pool

import (
	"github.com/derkan/nlog/common"
)

var _itemPool ItemPool

// Item is log item
type Item struct {
	callDepth  int
	subDepth   int
	formatters []common.Formatter
	lvl        common.Level
	pool       ItemPool
	buffs      []common.Buffer
	prefix     string
}

// NullItem is used for dismissed log levels
var NullItem = &Item{}

// Msg logs info message with format if logging level is satisfied
func (item *Item) Msg(args ...interface{}) {
	item.Msgf("", args...)
}

// Msgf logs info message with format if logging level is satisfied
func (item *Item) Msgf(format string, args ...interface{}) {
	if len(item.formatters) == 0 { // do no put back null item
		return
	}
	defer item.pool.put(item)
	for i := range item.formatters {
		if item.buffs != nil {
			item.formatters[i].Logf(item.callDepth+item.subDepth, item.lvl, item.buffs[i], item.prefix, format, args...)
		} else {
			item.formatters[i].Logf(item.callDepth+item.subDepth, item.lvl, nil, item.prefix, format, args...)
		}
	}
}

// Str adds a new str key value to buff
func (item *Item) Str(key string, val string) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Str(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Strs adds a slice of string value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Strs(key string, val []string) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Strs(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Int adds a new int key value to buff, Msg/Msgf should be called in same chain
func (item *Item) Int(key string, val int) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Int(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Ints adds a slice of int value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Ints(key string, val []int) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Ints(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Ints8 adds a slice of int8 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Ints8(key string, val []int8) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Ints8(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Ints16 adds a slice of int16 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Ints16(key string, val []int16) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Ints16(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Ints32 adds a slice of int32 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Ints32(key string, val []int32) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Ints32(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Int64 adds a new int64 key value to buff, Msg/Msgf should be called in same chain
func (item *Item) Int64(key string, val int64) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Int64(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Int64s adds a slice of int64 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Int64s(key string, val []int64) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Int64s(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// UInt adds a new uint key value to buff, Msg/Msgf should be called in same chain
func (item *Item) UInt(key string, val uint) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].UInt(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// UInts adds a slice of uint value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) UInts(key string, val []uint) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].UInts(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// UInts16 adds a slice of uint16 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) UInts16(key string, val []uint16) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].UInts16(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// UInts32 adds a slice of uint32 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) UInts32(key string, val []uint32) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].UInts32(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// UInt64 adds a new uint64 key value to buff, Msg/Msgf should be called in same chain
func (item *Item) UInt64(key string, val uint64) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].UInt64(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// UInts64 adds a slice of uint64 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) UInts64(key string, val []uint64) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].UInts64(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Float32 adds a new float32 key value to buff, Msg/Msgf should be called in same chain
func (item *Item) Float32(key string, val float32) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Float32(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Floats32 adds a slice of float32 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Floats32(key string, val []float32) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Floats32(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Float64 adds a new float64 key value to buff, Msg/Msgf should be called in same chain
func (item *Item) Float64(key string, val float64) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Float64(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Floats64 adds a slice of float32 value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Floats64(key string, val []float64) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Floats64(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Bool adds a new bool key value to buff, Msg/Msgf should be called in same chain
func (item *Item) Bool(key string, val bool) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Bool(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Bools adds a slice of bool value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Bools(key string, val []bool) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Bools(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// Error adds a new error key value to buff, Msg/Msgf should be called in same chain
func (item *Item) Err(val error) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}
	if val == nil {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Error(item.buffs[i], item.lvl, "err", val)
	}
	return item
}

// Bools adds a slice of error value with a key to buff, Msg/Msgf should be called in same chain
func (item *Item) Errors(key string, val []error) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}

	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].Errors(item.buffs[i], item.lvl, key, val)
	}
	return item
}

// With adds a new str key value to buff
func (item *Item) With(key string, val interface{}) common.LoggerItem {
	if len(item.formatters) == 0 {
		return item
	}
	// Each formatter will get a buffer, alloc space for Buffer objects
	if len(item.buffs) == 0 {
		item.buffs = make([]common.Buffer, len(item.formatters))
	}

	for i := range item.formatters {
		item.buffs[i] = GetBuffer()
		item.formatters[i].AppendKV(item.buffs[i], item.lvl, key, val)
	}
	return item
}
