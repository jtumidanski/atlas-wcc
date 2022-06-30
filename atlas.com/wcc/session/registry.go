package session

import (
	"sync"
)

type registry struct {
	mutex           sync.RWMutex
	sessionRegistry map[uint32]Model
	lockRegistry    map[uint32]*sync.RWMutex
}

var once sync.Once
var r *registry

func getRegistry() *registry {
	once.Do(func() {
		r = &registry{}
		r.sessionRegistry = make(map[uint32]Model)
		r.lockRegistry = make(map[uint32]*sync.RWMutex)
	})
	return r
}

func (r *registry) Add(s Model) {
	r.mutex.Lock()
	r.sessionRegistry[s.SessionId()] = s
	r.lockRegistry[s.SessionId()] = &sync.RWMutex{}
	r.mutex.Unlock()
}

func (r *registry) Remove(sessionId uint32) {
	r.mutex.Lock()
	delete(r.sessionRegistry, sessionId)
	delete(r.lockRegistry, sessionId)
	r.mutex.Unlock()
}

func (r *registry) Get(sessionId uint32) (Model, bool) {
	r.mutex.RLock()
	if s, ok := r.sessionRegistry[sessionId]; ok {
		r.mutex.RUnlock()
		return s, true
	}
	r.mutex.RUnlock()
	return Model{}, false
}

func (r *registry) GetLock(sessionId uint32) (*sync.RWMutex, bool) {
	r.mutex.RLock()
	if val, ok := r.lockRegistry[sessionId]; ok {
		r.mutex.RUnlock()
		return val, true
	}
	r.mutex.RUnlock()
	return nil, false
}

func (r *registry) GetAll() []Model {
	r.mutex.RLock()
	s := make([]Model, 0)
	for _, v := range r.sessionRegistry {
		s = append(s, v)
	}
	r.mutex.RUnlock()
	return s
}

func (r *registry) Update(m Model) {
	r.mutex.Lock()
	if _, ok := r.sessionRegistry[m.SessionId()]; ok {
		r.sessionRegistry[m.SessionId()] = m
	}
	r.mutex.Unlock()
}
