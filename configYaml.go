package caretakerd

import (
	"fmt"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
	"gopkg.in/blaubaer/yaml.v66"
	"io/ioutil"
	"os"
)

// ConfigDoesNotExistError descripts an error if a config does not exists.
type ConfigDoesNotExistError struct {
	fileName string
}

// Error returns the error message.
func (instance ConfigDoesNotExistError) Error() string {
	return fmt.Sprintf("Config '%v' does not exist.", instance.fileName)
}

// LoadFromYamlFile loads the caretakerd config from the given yaml file.
func LoadFromYamlFile(fileName values.String) (Config, error) {
	result := NewConfig()
	content, err := ioutil.ReadFile(fileName.String())
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ConfigDoesNotExistError{fileName: fileName.String()}
		}
		return Config{}, errors.New("Could not read config from '%v'.", fileName).CausedBy(err)
	}
	if err := yaml.Unmarshal(content, &result); err != nil {
		return Config{}, errors.New("Could not unmarshal config from '%v'.", fileName).CausedBy(err)
	}
	return result, nil
}

// WriteToYamlFile writes the config of the current instance to the given yaml file.
func (instance Config) WriteToYamlFile(fileName values.String) error {
	content, err := yaml.Marshal(instance)
	if err != nil {
		return errors.New("Could not write config to '%v'.", fileName).CausedBy(err)
	}
	if err := ioutil.WriteFile(fileName.String(), content, 0744); err != nil {
		return errors.New("Could not write marshalled config to '%v'.", fileName).CausedBy(err)
	}
	return nil
}
