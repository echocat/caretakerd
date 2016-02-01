package service

import (
	"github.com/echocat/caretakerd/errors"
	"strings"
)

// @inline
type Environments map[string]string

func (i Environments) String() string {
	result := ""
	for key, value := range i {
		if len(result) > 0 {
			result += "\n"
		}
		result += key + "=" + value
	}
	return result
}

func evaluate(value string) (map[string]string, error) {
	result := map[string]string{}
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return result, errors.New("Illegal environment settings format: %s", line)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

func (i *Environments) Set(value string) error {
	values, err := evaluate(value)
	if err != nil {
		return err
	}
	*i = values
	return nil
}

func (i Environments) Append(value string) error {
	values, err := evaluate(value)
	if err != nil {
		return err
	}
	for key, value := range values {
		i[key] = value
	}
	return nil
}

func (i *Environments) Put(key string, value string) error {
	(*i)[key] = value
	return nil
}
