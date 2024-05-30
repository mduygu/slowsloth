package statusprinter

import (
	"runtime"
	"sync/atomic"
)

type StatusManager struct {
	activeConnections int32
	serviceAvailable  int32
	totalRAMUsage     uint64
	totalBandwidth    uint64
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

func (sm *StatusManager) TotalRAMUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	totalRAMUsage := m.Sys / 1024 / 1024
	return totalRAMUsage
}

func (sm *StatusManager) IncrementTotalBandwidth(bytes uint64) {
	atomic.AddUint64(&sm.totalBandwidth, bytes)
}

func (sm *StatusManager) TotalBandwidth() uint64 {
	return atomic.LoadUint64(&sm.totalBandwidth)
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

func (sm *StatusManager) SetServiceColor(isAvailable bool) string {
	var color string

	if isAvailable {
		color = "\033[1m\033[32m"
	} else {
		color = "\033[1m\033[31m"
	}

	return color
}

func (sm *StatusManager) IsServiceAvailable() bool {
	return atomic.LoadInt32(&sm.serviceAvailable) == 1
}

func (sm *StatusManager) ServiceAvailability() string {

	serviceAvailability := "NO"
	if sm.IsServiceAvailable() {
		serviceAvailability = "YES"
	}
	return serviceAvailability
}
