package config

type Config struct {
	Selectors map[string]string `yaml:"selectors"`
	Ignore    []string          `yaml:"ignore"`
	Pins      map[string]string `yaml:"pins"`
}
