package manager

import (
	"sync"
)

var DefaultManager *Manager

type Manager struct {
	connections sync.Map
}

func (m *Manager) Join(c tracker) {
	m.connections.Store(c.ID(), c)
}

func (m *Manager) Leave(c tracker) {
	m.connections.Delete(c.ID())
}
