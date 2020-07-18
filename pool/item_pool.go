package pool

import (
	"sync"

	"github.com/derkan/nlog"
)

// A ItemPool is a type-safe wrapper around a sync.Pool.
type ItemPool struct {
	p *sync.Pool
}

// NewItemPool constructs a new ItemPool.
func NewItemPool(callDepth int, level nlog.Level, formatters ...nlog.Formatter) ItemPool {
	return ItemPool{p: &sync.Pool{
		New: func() interface{} {
			return &Item{callDepth: callDepth, lvl: level, formatters: formatters}
		},
	}}
}

// Get retrieves a Buffer from the pool
func (p ItemPool) Get(lvl nlog.Level, prefix string, subDepth int) *Item {
	item := p.p.Get().(*Item)
	item.pool = p
	item.lvl = lvl
	item.prefix = prefix
	item.subDepth = subDepth
	return item
}

func (p ItemPool) put(item *Item) {
	for i := range item.buffs {
		defer PutBuffer(item.buffs[i])
	}
	p.p.Put(item)
}
