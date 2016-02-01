package values

import (
	"strings"
)

// @inline
type ExitCodes []ExitCode

func (i ExitCodes) String() string {
	result := ""
	for _, code := range i {
		if len(result) > 0 {
			result += ","
		}
		result += code.String()
	}
	return result
}

func (i *ExitCodes) Set(value string) error {
	candidates := strings.Split(value, ",")
	result := ExitCodes{}
	for _, plainCandidate := range candidates {
		candidate := strings.TrimSpace(plainCandidate)
		if len(candidate) > 0 {
			var code ExitCode
			if err := code.Set(candidate); err != nil {
				return err
			}
			result = append(result, code)
		}
	}
	(*i) = result
	return nil
}

func (i ExitCodes) Validate() {
	i.String()
}

func (i ExitCodes) Contains(what ExitCode) bool {
	for _, candidate := range i {
		if candidate == what {
			return true
		}
	}
	return false
}
