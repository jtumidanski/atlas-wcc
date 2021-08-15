package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (c *Registry) load() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	con := &Model{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		return err
	}
	c.c = con
	return nil
}
