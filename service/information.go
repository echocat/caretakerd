package service

import (
	"github.com/echocat/caretakerd/values"
)

// Information represents the current status of a running execution.
type Information struct {
	Config Config         `json:"config"`
	Status Status         `json:"status"`
	PID    values.Integer `json:"pid"`
}

// NewInformationForExecution creates a new information instance for the given execution.
func NewInformationForExecution(e *Execution) Information {
	return Information{
		Config: e.service.config,
		Status: e.status,
		PID:    values.Integer(e.PID()),
	}
}

// NewInformationForService creates a new information instance for the given service.
// This always means that there is no execution and the service is currently down.
func NewInformationForService(s *Service) Information {
	return Information{
		Config: s.config,
		Status: Down,
		PID:    0,
	}
}
