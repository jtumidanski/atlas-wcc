package configuration

import (
	"sync"
)

type Registry struct {
	c *Model
	e error
}

var configurationRegistryOnce sync.Once
var configurationRegistry *Registry

func GetConfiguration() (*Model, error) {
	configurationRegistryOnce.Do(func() {
		configurationRegistry = &Registry{}
		err := configurationRegistry.loadConfiguration()
		configurationRegistry.e = err
	})
	return configurationRegistry.c, configurationRegistry.e
}
