package config

type Config struct {
	Selectors       map[string]string `json:"selectors" yaml:"selectors"`
	Allowed         []string          `json:"allowed" yaml:"allowed"`
	RefreshInterval int               `json:"refreshInterval" yaml:"refreshInterval"`
}
