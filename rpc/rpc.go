package rpc

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"log"
	"net"
	"strings"
	"strconv"
	"crypto/tls"
	"crypto/x509"
	. "github.com/echocat/caretakerd/values"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/control"
)

type Caretakerd interface {
	Control() *control.Control
	Services() *service.Services
	KeyStore() *keyStore.KeyStore
	Logger() *logger.Logger
	ConfigObject() interface{}
}

type Execution interface {
	GetFor(s *service.Service) (*service.Execution, bool)
	Information() map[string]service.Information
	InformationFor(s *service.Service) service.Information
	Start(*service.Service) error
	Restart(*service.Service) error
	Stop(*service.Service) error
	Kill(*service.Service) error
	Signal(*service.Service, Signal) error
}

type ListenerStopped struct{}

func (i ListenerStopped) Error() string {
	return "stopped"
}

type StoppableListener struct {
	*net.TCPListener
	stop chan int
}

func NewStoppableListener(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.TCPListener)
	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}
	result := &StoppableListener{
		TCPListener: tcpL,
		stop: make(chan int),
	}
	return result, nil
}

func (sl *StoppableListener) Accept() (net.Conn, error) {
	for {
		newConn, err := sl.TCPListener.Accept()
		if isClosedError(err) {
			return nil, ListenerStopped{}
		}
		return newConn, err
	}
}

func isClosedError(what error) bool {
	if what == nil {
		return false
	} else if opErr, ok := what.(*net.OpError); ok {
		message := strings.ToLower(opErr.Err.Error())
		return strings.Contains(message, "closed network connection")
	} else {
		return false
	}
}

type Rpc struct {
	conf       Config
	execution  Execution
	caretakerd Caretakerd
	listener   *StoppableListener
	logger     *logger.Logger
}

func NewRpc(conf Config, execution Execution, executable Caretakerd, log *logger.Logger) *Rpc {
	rpc := Rpc{
		conf: conf,
		execution: execution,
		caretakerd: executable,
		logger: log,
	}
	return &rpc
}

func (instance *Rpc) Start() {
	go instance.Run()
}

func (instance *Rpc) Run() {
	defer panics.DefaultPanicHandler()
	container := restful.NewContainer()

	ws := new(restful.WebService)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/config").To(instance.config))

	ws.Route(ws.GET("/control/config").To(instance.controlConfig))

	ws.Route(ws.GET("/services").To(instance.services))

	ws.Route(ws.GET("/service/{serviceName}").To(instance.service))
	ws.Route(ws.GET("/service/{serviceName}/config").To(instance.serviceConfig))
	ws.Route(ws.GET("/service/{serviceName}/state").To(instance.serviceStatus))
	ws.Route(ws.GET("/service/{serviceName}/pid").To(instance.servicePid))

	ws.Route(ws.POST("/service/{serviceName}/restart").To(instance.serviceRestart))
	ws.Route(ws.POST("/service/{serviceName}/status").To(instance.serviceStart))
	ws.Route(ws.POST("/service/{serviceName}/stop").To(instance.serviceStop))
	ws.Route(ws.POST("/service/{serviceName}/kill").To(instance.serviceKill))
	ws.Route(ws.POST("/service/{serviceName}/signal").To(instance.serviceSignal))

	container.Add(ws)

	server := &http.Server{
		Handler: container,
		ErrorLog: log.New(instance.logger.ReceiverFor(logger.Debug), "", 0),
	}
	instance.logger.Log(logger.Debug, "Rpc will bind to %v...", instance.conf.Listen)
	listener, err := net.Listen(instance.conf.Listen.AsScheme(), instance.conf.Listen.AsAddress())
	if err != nil {
		log.Fatal(err)
	}
	sl, err2 := NewStoppableListener(listener)
	if err2 != nil {
		panics.New("Could not create listener.").CausedBy(err2).Throw()
	}
	defer func() {
		(*instance).listener = nil
	}()
	(*instance).listener = sl
	if err := server.Serve(instance.secure(sl)); err != nil {
		if _, ok := err.(ListenerStopped); !ok {
			panics.New("Could not listen.").CausedBy(err2).Throw()
		}
	}
}

func (instance *Rpc) secure(in net.Listener) net.Listener {
	out := in
	sec := instance.caretakerd.KeyStore()
	keyPair, err := tls.X509KeyPair(sec.Pem(), sec.Pem())
	if err != nil {
		panics.New("Could not load pem of caretakerd.").CausedBy(err).Throw()
	}

	rootCas := x509.NewCertPool()
	for _, cert := range sec.Ca() {
		rootCas.AddCert(cert)
	}

	out = tls.NewListener(in, &tls.Config{
		NextProtos: []string{"http/1.1"},
		Certificates: []tls.Certificate{keyPair},
		RootCAs: rootCas,
		ClientCAs: rootCas,
		ClientAuth: tls.RequireAndVerifyClientCert,
	})

	return out
}

func (instance *Rpc) Stop() {
	listener := (*instance).listener
	if listener != nil {
		listener.Close()
	}
}

func (instance *Rpc) checkPermission(request *restful.Request, permissionChecker func(access.Access) bool) bool {
	if request != nil {
		hr := request.Request
		if hr != nil {
			cs := hr.TLS
			if cs != nil {
				for _, cert := range cs.PeerCertificates {
					ctl := instance.caretakerd.Control()
					acc := ctl.Access()
					if acc.IsCertValid(cert) && permissionChecker(*acc) {
						return true
					}
					for _, serv := range (*instance.caretakerd.Services()) {
						acc := serv.Access()
						if acc.IsCertValid(cert) && permissionChecker(*acc) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (instance *Rpc) hasReadPermission(request *restful.Request) bool {
	return instance.checkPermission(request, func(a access.Access) bool {
		return a.HasReadPermission()
	})
}

func (instance *Rpc) hasWritePermission(request *restful.Request) bool {
	return instance.checkPermission(request, func(a access.Access) bool {
		return a.HasWritePermission()
	})
}

func (instance *Rpc) onReadPermission(request *restful.Request, response *restful.Response, doThis func()) {
	if instance.hasReadPermission(request) {
		doThis()
	} else {
		response.WriteErrorString(http.StatusForbidden, "No read permission to instance endpoint.")
	}
}

func (instance *Rpc) onWritePermission(request *restful.Request, response *restful.Response, doThis func()) {
	if instance.hasWritePermission(request) {
		doThis()
	} else {
		response.WriteErrorString(http.StatusForbidden, "No write permission to instance endpoint.")
	}
}

func (instance *Rpc) config(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		response.WriteEntity(instance.caretakerd.ConfigObject())
	})
}

func (instance *Rpc) controlConfig(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		response.WriteEntity(instance.caretakerd.Control().ConfigObject())
	})
}

func (instance *Rpc) services(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		information := instance.execution.Information()
		response.WriteEntity(information)
	})
}

func (instance *Rpc) service(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		serviceName := request.PathParameter("serviceName")
		services := instance.caretakerd.Services()
		if service, ok := services.Get(serviceName); ok {
			information := instance.execution.InformationFor(service)
			response.WriteEntity(information)
		} else {
			response.WriteError(http.StatusNotFound, errors.New("Service '%s' does not exist.", serviceName))
		}
	})
}

func (instance *Rpc) serviceConfig(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		instance.doWithService(request, response, func(service *service.Service) {
			response.WriteEntity(service.Config())
		})
	})
}

func (instance *Rpc) serviceStatus(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		instance.doWithExecution(request, response, func(execution *service.Execution) {
			if execution != nil {
				response.Write([]byte(execution.Status().String()))
			} else {
				response.Write([]byte(service.Down.String()))
			}
		})
	})
}

func (instance *Rpc) servicePid(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		instance.doWithExecution(request, response, func(execution *service.Execution) {
			if execution != nil {
				response.Write([]byte(strconv.Itoa(execution.Pid())))
			} else {
				response.Write([]byte("0"))
			}
		})
	})
}

func (instance *Rpc) serviceRestart(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(service *service.Service) {
			err := instance.execution.Restart(service)
			if err == nil {
				response.Write([]byte("OK"))
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: " + err.Error())
			}
		})
	})
}

func (instance *Rpc) serviceStart(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			err := instance.execution.Start(sc)
			if err == nil {
				response.Write([]byte("OK"))
			} else if sde, ok := err.(service.ServiceAlreadyRunningError); ok {
				response.WriteErrorString(http.StatusConflict, "ERROR: " + sde.Error())
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: " + err.Error())
			}
		})
	})
}

func (instance *Rpc) serviceStop(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			err := instance.execution.Stop(sc)
			if err == nil {
				response.Write([]byte("OK"))
			} else if sde, ok := err.(service.ServiceDownError); ok {
				response.WriteErrorString(http.StatusConflict, "ERROR: " + sde.Error())
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: " + err.Error())
			}
		})
	})
}

func (instance *Rpc) serviceKill(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			err := instance.execution.Kill(sc)
			if err == nil {
				response.Write([]byte("OK"))
			} else if sde, ok := err.(service.ServiceDownError); ok {
				response.WriteErrorString(http.StatusConflict, "ERROR: " + sde.Error())
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: " + err.Error())
			}
		})
	})
}

type SignalBody struct {
	Signal Signal `json:"signal"`
}

func (instance *Rpc) serviceSignal(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			sb := SignalBody{}
			err := request.ReadEntity(&sb)
			if err != nil {
				response.WriteErrorString(http.StatusBadRequest, "ERROR: Illegal body. " + err.Error())
			} else {
				err = instance.execution.Kill(sc)
				if err == nil {
					response.Write([]byte("OK"))
				} else if sde, ok := err.(service.ServiceDownError); ok {
					response.WriteErrorString(http.StatusConflict, "ERROR: " + sde.Error())
				} else {
					response.WriteErrorString(http.StatusInternalServerError, "ERROR: " + err.Error())
				}
			}
		})
	})
}

func (instance *Rpc) doWithService(request *restful.Request, response *restful.Response, what func(service *service.Service)) {
	serviceName := request.PathParameter("serviceName")
	services := instance.caretakerd.Services()
	if service, ok := services.Get(serviceName); ok {
		what(service)
	} else {
		response.WriteError(http.StatusNotFound, errors.New("Service '%s' does not exist.", serviceName))
	}
}

func (instance *Rpc) doWithExecution(request *restful.Request, response *restful.Response, what func(execution *service.Execution)) {
	instance.doWithService(request, response, func(s *service.Service) {
		if e, ok := instance.execution.GetFor(s); ok {
			what(e)
		} else {
			what(nil)
		}
	})
}

