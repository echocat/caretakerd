package service

import (
    . "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/service/status"
    "github.com/echocat/caretakerd/service/config"
)

type Information struct {
    Config config.Config `json:"config"`
    Status status.Status `json:"status"`
    Pid    Integer `json:"pid"`
}

func NewInformationForExecution(e *Execution) Information {
    return Information{
        Config: e.service.config,
        Status: e.status,
        Pid: Integer(e.Pid()),
    }
}

func NewInformationForService(s *Service) Information {
    return Information{
        Config: s.config,
        Status: status.Down,
        Pid: 0,
    }
}
