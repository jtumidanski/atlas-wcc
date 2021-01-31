package registries

import (
	"atlas-wcc/mapleSession"
	"sync"
)

type SessionRegistry struct {
	mutex           sync.RWMutex
	sessionRegistry map[int]*mapleSession.MapleSession
}

var sessionRegistryOnce sync.Once
var sessionRegistry *SessionRegistry

func GetSessionRegistry() *SessionRegistry {
	sessionRegistryOnce.Do(func() {
		sessionRegistry = &SessionRegistry{}
		sessionRegistry.sessionRegistry = make(map[int]*mapleSession.MapleSession)
	})
	return sessionRegistry
}

func (r *SessionRegistry) Add(s *mapleSession.MapleSession) {
	r.mutex.Lock()
	r.sessionRegistry[(*s).SessionId()] = s
	r.mutex.Unlock()
}

func (r *SessionRegistry) Remove(sessionId int) {
	r.mutex.Lock()
	delete(r.sessionRegistry, sessionId)
	r.mutex.Unlock()
}

func (r *SessionRegistry) Get(sessionId int) mapleSession.MapleSession {
	r.mutex.RLock()
	s := r.sessionRegistry[sessionId]
	r.mutex.RUnlock()
	return *s
}

func (r *SessionRegistry) GetAll() []mapleSession.MapleSession {
	r.mutex.RLock()
	s := make([]mapleSession.MapleSession, 0)
	for _, v := range r.sessionRegistry {
		s = append(s, *v)
	}
	r.mutex.RUnlock()
	return s
}
