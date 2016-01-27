package service

type ServiceDownError struct {
	Name string
}

func (instance ServiceDownError) Error() string {
	return "Service '" + instance.Name + "' is down."
}

type ServiceAlreadyRunningError struct {
	Name string
}

func (instance ServiceAlreadyRunningError) Error() string {
	return "Service '" + instance.Name + "' already running."
}

