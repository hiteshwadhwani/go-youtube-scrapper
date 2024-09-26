package apikeymanager

import "sync"

type Manager struct {
	keys          []string
	currentKey    int
	quotaExceeded []bool
	mutex         sync.Mutex
}

func New(keys []string) *Manager {
	return &Manager{
		keys:          keys,
		currentKey:    0,
		quotaExceeded: make([]bool, len(keys)),
	}
}

// get next key that is not exceeded quota
func (m *Manager) GetNextKey() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i := range m.keys {
		currentidx := (m.currentKey + i) % len(m.keys)
		if !m.quotaExceeded[currentidx] {
			m.currentKey = currentidx
			return m.keys[currentidx]
		}
	}

	return ""
}

// mark current key as exceeded quota
func (m *Manager) MarkQuotaExceed() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.quotaExceeded[m.currentKey] = true
}
