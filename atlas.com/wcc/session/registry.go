package session

import (
	"sync"
)

type Registry struct {
	mutex           sync.RWMutex
	sessionRegistry map[uint32]*Model
}

var sessionRegistryOnce sync.Once
var sessionRegistry *Registry

func GetRegistry() *Registry {
	sessionRegistryOnce.Do(func() {
		sessionRegistry = &Registry{}
		sessionRegistry.sessionRegistry = make(map[uint32]*Model)
	})
	return sessionRegistry
}

func (r *Registry) Add(s *Model) {
	r.mutex.Lock()
	r.sessionRegistry[s.SessionId()] = s
	r.mutex.Unlock()
}

func (r *Registry) Remove(sessionId uint32) {
	r.mutex.Lock()
	delete(r.sessionRegistry, sessionId)
	r.mutex.Unlock()
}

func (r *Registry) Get(sessionId uint32) *Model {
	r.mutex.RLock()
	s := r.sessionRegistry[sessionId]
	r.mutex.RUnlock()
	return s
}

func (r *Registry) GetAll() []*Model {
	r.mutex.RLock()
	s := make([]*Model, 0)
	for _, v := range r.sessionRegistry {
		s = append(s, v)
	}
	r.mutex.RUnlock()
	return s
}
