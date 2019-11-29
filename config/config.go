package config

import (
	"github.com/loganprice/untitled/config/env"
	"github.com/loganprice/untitled/config/flag"
)

type (
	// Sources stores the ConfServices that will be used
	Sources map[string]ConfService
	// ConfService provides operations to read and unmarshal config
	ConfService interface {
		GetValues() error
		Unmarshal(interface{}) error
	}
)

// NewSourceStore Creates an empty Sources Object
func NewSourceStore() Sources {
	return Sources{}
}

// NewDefaultSourceStore Creates a Sources Object that contains config sources of flags and environment
func NewDefaultSourceStore() Sources {
	return Sources{
		"args-default": flag.NewConfig(),
		"env-default":  env.NewConfig(),
	}
}

// AddSource gives the ability to add a new source to the default list of registered Sources
func (s Sources) AddSource(sourceName string, source ConfService) Sources {
	s[sourceName] = source
	return s
}

// Load takes in arguments of the source that you want to pull config from and the object you want to
// unmarshal that data into. There is a special source called "merge" which will unmarshal all your registered
// source in alphabetical order. It is recommend that you name you prefix based on the precedence (Command Line Args, File,  Env,  External, nonTyped)
func (s Sources) Load(source string, targets ...interface{}) error {
	for _, target := range targets {
		if source == "merge" {
			orderedSources := setPrecedence(s)
			for _, sourceKey := range orderedSources {
				if err := getAndUnmarshal(s[sourceKey], target); err != nil {
					return err
				}
			}
			return nil
		}
		if err := getAndUnmarshal(s[source], target); err != nil {
			return err
		}
	}
	return nil
}

func getAndUnmarshal(c ConfService, target interface{}) error {
	if err := c.GetValues(); err != nil {
		return err
	}
	return c.Unmarshal(target)
}
