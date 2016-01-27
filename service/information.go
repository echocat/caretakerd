package service

import (
	. "github.com/echocat/caretakerd/values"
)

type Information struct {
	Config Config  `json:"config"`
	Status Status  `json:"status"`
	Pid    Integer `json:"pid"`
}

func NewInformationForExecution(e *Execution) Information {
	return Information{
		Config: e.service.config,
		Status: e.status,
		Pid:    Integer(e.Pid()),
	}
}

func NewInformationForService(s *Service) Information {
	return Information{
		Config: s.config,
		Status: Down,
		Pid:    0,
	}
}
