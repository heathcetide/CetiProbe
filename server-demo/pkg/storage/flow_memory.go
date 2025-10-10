package storage

import (
	"sync"
	"time"

	"probe/internal/models"
)

type FlowStorage interface {
	Add(flow *models.Flow)
	GetAll(limit int) []*models.Flow
	GetByID(id string) *models.Flow
	Clear()
	Stats() FlowStats
}

type FlowStats struct {
	Total     int       `json:"total"`
	StartTime time.Time `json:"start_time"`
	LastTime  time.Time `json:"last_time"`
}

type memoryFlowStore struct {
	mu    sync.RWMutex
	flows []*models.Flow
	byID  map[string]*models.Flow
	stats FlowStats
}

func NewMemoryFlowStore() FlowStorage {
	return &memoryFlowStore{
		flows: make([]*models.Flow, 0, 1024),
		byID:  make(map[string]*models.Flow),
		stats: FlowStats{StartTime: time.Now()},
	}
}

func (m *memoryFlowStore) Add(flow *models.Flow) {
	m.mu.Lock()
	defer m.mu.Unlock()
	const max = 20000
	if len(m.flows) >= max {
		old := m.flows[0]
		delete(m.byID, old.ID)
		m.flows = m.flows[1:]
	}
	m.flows = append(m.flows, flow)
	m.byID[flow.ID] = flow
	m.stats.Total++
	m.stats.LastTime = time.Now()
}

func (m *memoryFlowStore) GetAll(limit int) []*models.Flow {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if limit <= 0 || limit > len(m.flows) {
		limit = len(m.flows)
	}
	start := len(m.flows) - limit
	return m.flows[start:]
}

func (m *memoryFlowStore) GetByID(id string) *models.Flow {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byID[id]
}

func (m *memoryFlowStore) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flows = make([]*models.Flow, 0, 1024)
	m.byID = make(map[string]*models.Flow)
	m.stats = FlowStats{StartTime: time.Now()}
}

func (m *memoryFlowStore) Stats() FlowStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stats
}
