package service

// AlreadyRunningError indicates that a service was up but was expected to be down.
type AlreadyRunningError struct {
	Name string
}

func (instance AlreadyRunningError) Error() string {
	return "Service '" + instance.Name + "' already running."
}

// AlreadyStoppedError indicates that a service was up but was expected to be down.
type AlreadyStoppedError struct {
	Name string
}

func (instance AlreadyStoppedError) Error() string {
	return "Service '" + instance.Name + "' already stopped."
}
