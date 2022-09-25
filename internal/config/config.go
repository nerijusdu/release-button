package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ReadConfig() (*Config, error) {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(file, c)
	return c, err
}

func WriteConfig(c Config) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.yaml", bytes, 0644)
	return err
}
