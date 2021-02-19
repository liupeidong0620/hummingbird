package nat

import (
	"net"
	"sync"
	"sync/atomic"
)

type Entry struct {
	count   int32
	mapping sync.Map
}

type Table struct {
	mapping sync.Map
}

func (e *Entry) Set(key string, conn net.Conn) {
	atomic.AddInt32(&e.count, 1)
	e.mapping.Store(key, conn)
}

func (e *Entry) Get(key string) net.Conn {
	item, exist := e.mapping.Load(key)
	if !exist {
		return nil
	}
	return item.(net.Conn)
}

func (e *Entry) Delete(key string) {
	atomic.AddInt32(&e.count, -1)
	e.mapping.Delete(key)
}

func (e *Entry) IsEmpty() bool {
	if atomic.LoadInt32(&e.count) <= 0 {
		return true
	}

	return false
}

func (t *Table) Set(key string, entry *Entry) {
	t.mapping.Store(key, entry)
}

func (t *Table) Get(key string) *Entry {
	item, exist := t.mapping.Load(key)
	if !exist {
		return nil
	}
	return item.(*Entry)
}

func (t *Table) GetOrCreateLock(key string) (*sync.Cond, bool) {
	item, loaded := t.mapping.LoadOrStore(key, sync.NewCond(&sync.Mutex{}))
	return item.(*sync.Cond), loaded
}

func (t *Table) Delete(key string) {
	t.mapping.Delete(key)
}

func NewTable() *Table {
	return &Table{}
}

func NewEntry() *Entry {
	return &Entry{}
}
