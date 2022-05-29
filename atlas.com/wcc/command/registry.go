package command

import (
	"atlas-wcc/session"
	"sync"
)

type registry struct {
	commandRegistry []Producer
}

var once sync.Once
var r *registry

func Registry() *registry {
	once.Do(func() {
		r = &registry{}
		r.commandRegistry = make([]Producer, 0)
	})
	return r
}

func (r *registry) Add(svs ...Producer) {
	for _, sv := range svs {
		r.commandRegistry = append(r.commandRegistry, sv)
	}
}

func (r *registry) Get(s session.Model, m string) (Executor, bool) {
	for _, c := range r.commandRegistry {
		e, found := c(s, m)
		if found {
			return e, found
		}
	}
	return nil, false
}
