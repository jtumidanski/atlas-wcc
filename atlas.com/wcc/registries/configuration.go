package registries

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type ConfigurationRegistry struct {
	c *Configuration
	e error
}

type Configuration struct {
	TimeoutTaskInterval int64 `yaml:"timeoutTaskInterval"`
	TimeoutDuration     int64 `yaml:"timeoutDuration"`
}

var configurationRegistryOnce sync.Once
var configurationRegistry *ConfigurationRegistry

func GetConfiguration() (*Configuration, error) {
	configurationRegistryOnce.Do(func() {
		configurationRegistry = &ConfigurationRegistry{}
		err := configurationRegistry.loadConfiguration()
		configurationRegistry.e = err
	})
	return configurationRegistry.c, configurationRegistry.e
}

func (c *ConfigurationRegistry) loadConfiguration() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		return err
	}
	c.c = con
	return nil
}
