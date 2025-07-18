package importer

import "sync"

type RefreshManager struct {
	mu      sync.Mutex
	running bool
}

func NewManager() *RefreshManager {
	return &RefreshManager{}
}

func (r *RefreshManager) TryLock() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return false
	} else {
		r.running = true
		return true
	}
}

func (r *RefreshManager) Unlock() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.running = false
}

func (r *RefreshManager) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}
