package caretakerd

import (
	"io/ioutil"
	. "github.com/echocat/caretakerd/values"
	"gopkg.in/blaubaer/goyaml.v66"
	"github.com/echocat/caretakerd/errors"
	"os"
	"fmt"
)

type ConfigDoesNotExistError struct {
	fileName string
}

func (instance ConfigDoesNotExistError) Error() string {
	return fmt.Sprintf("Config '%v' does not exist.", instance.fileName)
}

func LoadFromYamlFile(fileName String) (Config, error) {
	result := NewConfig()
	content, err := ioutil.ReadFile(fileName.String())
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ConfigDoesNotExistError{fileName: fileName.String()}
		} else {
			return Config{}, errors.New("Could not read config from '%v'.", fileName).CausedBy(err)
		}
	}
	if err := yaml.Unmarshal(content, &result); err != nil {
		return Config{}, errors.New("Could not unmarshal config from '%v'.", fileName).CausedBy(err)
	}
	return result, nil
}

func (i Config) WriteToYamlFile(fileName String) error {
	content, err := yaml.Marshal(i)
	if err != nil {
		return errors.New("Could not write config to '%v'.", fileName).CausedBy(err)
	}
	if err := ioutil.WriteFile(fileName.String(), content, 0744); err != nil {
		return errors.New("Could not write marshalled config to '%v'.", fileName).CausedBy(err)
	}
	return nil
}
