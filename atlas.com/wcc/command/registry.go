package command

import (
	"atlas-wcc/session"
	"sync"
)

type registry struct {
	commandRegistry []*model
}

var once sync.Once
var r *registry

type model struct {
	sv SyntaxValidator
	e  Executor
}

func (m *model) Executor() Executor {
	return m.e
}

func (m *model) SyntaxValidator() SyntaxValidator {
	return m.sv
}

func Registry() *registry {
	once.Do(func() {
		r = &registry{}
		r.commandRegistry = make([]*model, 0)
	})
	return r
}

func (r *registry) Add(sv SyntaxValidator, e Executor) {
	r.commandRegistry = append(r.commandRegistry, &model{
		sv: sv,
		e:  e,
	})
}

func (r *registry) Get(s session.Model, m string) (Executor, bool) {
	for _, c := range r.commandRegistry {
		if c.SyntaxValidator()(s, m) {
			return c.Executor(), true
		}
	}
	return nil, false
}
