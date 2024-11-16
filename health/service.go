package health

import "sync/atomic"

// Service health status
type Health struct {
	status int32 // atomic
}

func (h *Health) SetHealthy() {
	atomic.StoreInt32(&h.status, 1)
}

func (h *Health) SetUnhealthy() {
	atomic.StoreInt32(&h.status, 0)
}

func (h *Health) IsHealthy() bool {
	return atomic.LoadInt32(&h.status) == 1
}
