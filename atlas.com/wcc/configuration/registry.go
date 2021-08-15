package configuration

import (
	"sync"
)

type Registry struct {
	c *Model
	e error
}

var once sync.Once
var registry *Registry

func Get() (*Model, error) {
	once.Do(func() {
		registry = &Registry{}
		err := registry.load()
		registry.e = err
	})
	return registry.c, registry.e
}
