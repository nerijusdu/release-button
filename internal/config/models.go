package config

type Config struct {
	Selectors       map[string]string `json:"selectors" yaml:"selectors"`
	Ignore          []string          `json:"ignore" yaml:"ignore"`
	Allowed         []string          `json:"allowed" yaml:"allowed"`
	RefreshInterval int               `json:"refreshInterval" yaml:"refreshInterval"`
}
