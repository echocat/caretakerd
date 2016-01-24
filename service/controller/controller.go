package controller
import (
    "sync"
    "github.com/echocat/caretakerd/service"
    . "github.com/echocat/caretakerd/service/exitCode"
    "github.com/echocat/caretakerd/service/signal"
)

type Controller struct {
    waitGroup *sync.WaitGroup
    services  *service.Services
}

func NewController(services service.Services) *Controller {
    return &Controller{
        waitGroup: new(sync.WaitGroup),
        services: services,
    }
}

func (i *Controller) Run() (ExitCode, error) {
    notificationHandler := InstallNotificationHandler(i.handleTerminationSignal)
    defer i.finalizeRun(notificationHandler)

}

func (i *Controller) startService(service *service.Service) (ExitCode, error) {
    return service.Run()
}

func (i *Controller) runService(service *service.Service) (ExitCode, error) {
    return service.Run()
}

func (i *Controller) finalizeRun(notificationHandler *NotificationHandler) {
    notificationHandler.Uninstall()
}

func (i *Controller) handleTerminationSignal(signal signal.Signal) {
    i.StopAll()
}

func (i *Controller) Stop(service *service.Service) {
    service.Stop()
}

func (i *Controller) StopAll() {
    master := i.services.GetMaster()
    if master != nil {
        i.Stop(master)
    }
    for _, service := range i.services.GetAllButMaster() {
        i.Stop(service)
    }
}
