package statusprinter

import (
	"sync/atomic"
)

type StatusManager struct {
	activeConnections int32
	serviceAvailable  int32 // Using int32 to allow atomic operations
}

func NewStatusManager() *StatusManager {
	return &StatusManager{}
}

func (sm *StatusManager) IncrementActiveConnections() {
	atomic.AddInt32(&sm.activeConnections, 1)
}

func (sm *StatusManager) DecrementActiveConnections() {
	atomic.AddInt32(&sm.activeConnections, -1)
}

func (sm *StatusManager) ActiveConnections() int32 {
	return atomic.LoadInt32(&sm.activeConnections)
}

func (sm *StatusManager) SetServiceAvailable(isAvailable bool) {
	var val int32
	if isAvailable {
		val = 1
	} else {
		val = 0
	}
	atomic.StoreInt32(&sm.serviceAvailable, val)
}

func (sm *StatusManager) IsServiceAvailable() bool {
	return atomic.LoadInt32(&sm.serviceAvailable) == 1
}
