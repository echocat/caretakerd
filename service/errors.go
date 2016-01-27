package service

type ServiceDownError struct {
	Name string
}

func (this ServiceDownError) Error() string {
	return "Service '" + this.Name + "' is down."
}

type ServiceAlreadyRunningError struct {
	Name string
}

func (this ServiceAlreadyRunningError) Error() string {
	return "Service '" + this.Name + "' already running."
}

