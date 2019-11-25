package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Sources []string
	files   []file
	err     error
}

type file struct {
	Path  string
	Order int
	Data  []byte
}

// FromFile takes in an array of file locaations and inorder reads the files and unmarshals the bytes to
// the target. Files in the array
// func FromFile(target interface{}, fileLocations ...string) error {
// 	temp := &Values{}
// 	return temp.GetFileData(fileLocations...).UnmarshalFileData(target)

func NewConfig(fileLocations ...string) *Config {
	return &Config{
		Sources: fileLocations,
	}
}

// GetValues ... Builder
func (c *Config) GetValues() error {
	for i, fileLocation := range c.Sources {

		// Read in Contents
		fileContents, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			return err
		}
		c.files = append(c.files, file{
			Path:  fileLocation,
			Order: i,
			Data:  fileContents,
		})
	}
	return nil
}

// Unmarshal ...
func (c *Config) Unmarshal(target interface{}) error {
	// Sort Files to ensures that they are in the correct order
	sort.Slice(c.files, func(i, j int) bool {
		return c.files[i].Order < c.files[j].Order
	})
	for _, f := range c.files {
		switch filepath.Ext(f.Path) {
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(f.Data, target); err != nil {
				return err
			}
		case ".json":
			if err := json.Unmarshal(f.Data, target); err != nil {
				return err
			}
		case ".toml":
			if err := toml.Unmarshal(f.Data, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("couldn't determine file type of %v to unmarshal into Target", f.Path)
		}
	}
	return nil
}
