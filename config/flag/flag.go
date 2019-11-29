//Package flag - limitations cannot capture and array of flags
package flag

import (
	"os"
	"regexp"

	configutils "github.com/loganprice/untitled/config/utils"
)

const (
	flagWithEquals = `^-+(\S+=\S+)`
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
	flags := filterOutList(flagWithEquals, os.Args[1:])
	c.values = configutils.ParseKVStringSlice(flags, "=")
	return nil
}

// Unmarshal ...
func (c *Config) Unmarshal(target interface{}) error {
	err := configutils.Unmarshal(c.values, "short_flag", target)
	if err != nil {
		return err
	}
	return configutils.Unmarshal(c.values, "long_flag", target)
}

func filterOutList(permitted string, input []string) []string {
	var temp []string
	for _, value := range input {
		if re := regexp.MustCompile(permitted); re.MatchString(value) {
			temp = append(temp, re.FindStringSubmatch(value)[1])
		}
	}
	return temp
}
