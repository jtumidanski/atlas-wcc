package session

import (
   "sync"
)

type registry struct {
   mutex           sync.RWMutex
   sessionRegistry map[uint32]*Model
}

var once sync.Once
var r *registry

func Registry() *registry {
   once.Do(func() {
      r = &registry{}
      r.sessionRegistry = make(map[uint32]*Model)
   })
   return r
}

func (r *registry) Add(s *Model) {
   r.mutex.Lock()
   r.sessionRegistry[s.SessionId()] = s
   r.mutex.Unlock()
}

func (r *registry) Remove(sessionId uint32) {
   r.mutex.Lock()
   delete(r.sessionRegistry, sessionId)
   r.mutex.Unlock()
}

func (r *registry) Get(sessionId uint32) *Model {
   r.mutex.RLock()
   s := r.sessionRegistry[sessionId]
   r.mutex.RUnlock()
   return s
}

func (r *registry) GetAll() []*Model {
   r.mutex.RLock()
   s := make([]*Model, 0)
   for _, v := range r.sessionRegistry {
      s = append(s, v)
   }
   r.mutex.RUnlock()
   return s
}
