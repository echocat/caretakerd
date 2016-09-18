package values

import (
	"strings"
)

// ExitCodes represents a couple of exitCodes.
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

// Set sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
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

// Validate validates actions on this object and returns an error object if there are any.
func (i ExitCodes) Validate() {}

// Contains returns "true" if the given exitCode (what) is contained in this exitCodes.
func (i ExitCodes) Contains(what ExitCode) bool {
	for _, candidate := range i {
		if candidate == what {
			return true
		}
	}
	return false
}
