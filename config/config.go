package config

import (
	"github.com/loganprice/avalanche/config/env"
)

// ConfService provides operations to read and unmashal config
type ConfService interface {
	GetValues() error
	Unmarshal(interface{}) error
}

const (
	defaultFileName = "application.yml"
)

var registeredSources = map[string]ConfService{
	"defaultEnv": env.NewConfig(),
}

// RegisterSource gives the ability to add a new source to the default list of registered Sources
func RegisterSource(sourceName string, source ConfService) {
	registeredSources[sourceName] = source
}

// Load takes in arguments of the source that you want to pull config from and the object you want to
// unmarshal that data into. There is a special source called "merge" which will unmarshal all your registered
// source based on the persidence (File,  Env, Command Line Args, External)
func Load(source string, target interface{}) error {
	if source == "merge" {
		for _, source := range registeredSources {
			getAndUnmarshal(source, target)
		}
		return nil
	}
	return getAndUnmarshal(registeredSources[source], target)

}

func getAndUnmarshal(c ConfService, target interface{}) error {
	if err := c.GetValues(); err != nil {
		return err
	}
	return c.Unmarshal(target)
}

// To DO - Create a function to sort ConfServices based on persidence
// Load Priority - Default, File,  Env, Command Line Args, External