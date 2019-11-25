package env

import (
	configutils "github.com/loganprice/untitled/config/utils"
	"os"
)

// Config ...
type Config struct {
	values map[string]string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// GetValues ...
func (c *Config) GetValues() error {
	env := os.Environ()
	c.values = configutils.ParseKVStringSlice(env, "=")
	return nil
}

// Unmarshal ...
func (c *Config) Unmarshal(target interface{}) error {
	return configutils.Unmarshal(c.values, "env", target)
}
