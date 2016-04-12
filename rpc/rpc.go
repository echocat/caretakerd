package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/control"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/values"
	"github.com/emicklei/go-restful"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Caretakerd represents a caretakerd instance.
type Caretakerd interface {
	Control() *control.Control
	Services() *service.Services
	KeyStore() *keyStore.KeyStore
	Logger() *logger.Logger
	ConfigObject() interface{}
}

// Execution represents a caretakerd execution instance.
type Execution interface {
	GetFor(s *service.Service) (*service.Execution, bool)
	Information() map[string]service.Information
	InformationFor(s *service.Service) service.Information
	Start(*service.Service) error
	Restart(*service.Service) error
	Stop(*service.Service) error
	Kill(*service.Service) error
	Signal(*service.Service, values.Signal) error
}

// ListenerStoppedError occurs if the network listener is already stopped.
type ListenerStoppedError struct{}

func (i ListenerStoppedError) Error() string {
	return "stopped"
}

// StoppableListener is a reimplementation of net.TCPListener which is graceful stoppable.
type StoppableListener struct {
	*net.TCPListener
	stop chan int
}

// NewStoppableListener creates a new instance of StoppableListener and encapsulate the given Listener.
func NewStoppableListener(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.TCPListener)
	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}
	result := &StoppableListener{
		TCPListener: tcpL,
		stop:        make(chan int),
	}
	return result, nil
}

// Accept returns new connection if a remote client connects to the server.
// This method is blocking.
func (sl *StoppableListener) Accept() (net.Conn, error) {
	for {
		newConn, err := sl.TCPListener.Accept()
		if isClosedError(err) {
			return nil, ListenerStoppedError{}
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

// RPC holds all listeners and other resources of the RPC mechanism.
type RPC struct {
	conf       Config
	execution  Execution
	caretakerd Caretakerd
	listener   *StoppableListener
	logger     *logger.Logger
}

// NewRPC creates a new instance of RPC.
func NewRPC(conf Config, execution Execution, executable Caretakerd, log *logger.Logger) *RPC {
	rpc := RPC{
		conf:       conf,
		execution:  execution,
		caretakerd: executable,
		logger:     log,
	}
	return &rpc
}

// Start starts the RPC instance in background.
// This means: This method is not blocking.
func (instance *RPC) Start() {
	go instance.Run()
}

// Run starts the RPC instance in foreground.
// This means: This method is blocking.
func (instance *RPC) Run() {
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

	ws.Route(ws.POST("/service/{serviceName}/start").To(instance.serviceStart))
	ws.Route(ws.POST("/service/{serviceName}/restart").To(instance.serviceRestart))
	ws.Route(ws.POST("/service/{serviceName}/status").To(instance.serviceStart))
	ws.Route(ws.POST("/service/{serviceName}/stop").To(instance.serviceStop))
	ws.Route(ws.POST("/service/{serviceName}/kill").To(instance.serviceKill))
	ws.Route(ws.POST("/service/{serviceName}/signal").To(instance.serviceSignal))

	container.Add(ws)

	server := &http.Server{
		Handler:  container,
		ErrorLog: log.New(instance.logger.NewOutputStreamWrapperFor(logger.Debug), "", 0),
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
		if _, ok := err.(ListenerStoppedError); !ok {
			panics.New("Could not listen.").CausedBy(err2).Throw()
		}
	}
}

func (instance *RPC) secure(in net.Listener) net.Listener {
	out := in
	sec := instance.caretakerd.KeyStore()
	keyPair, err := tls.X509KeyPair(sec.PEM(), sec.PEM())
	if err != nil {
		panics.New("Could not load pem of caretakerd.").CausedBy(err).Throw()
	}

	rootCas := x509.NewCertPool()
	for _, cert := range sec.CA() {
		rootCas.AddCert(cert)
	}

	out = tls.NewListener(in, &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{keyPair},
		RootCAs:      rootCas,
		ClientCAs:    rootCas,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return out
}

// Stop stops the current RPC instance if running.
// This method is blocking.
func (instance *RPC) Stop() {
	listener := (*instance).listener
	if listener != nil {
		listener.Close()
	}
}

func (instance *RPC) checkPermission(request *restful.Request, permissionChecker func(access.Access) bool) bool {
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
					for _, serv := range *instance.caretakerd.Services() {
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

func (instance *RPC) hasReadPermission(request *restful.Request) bool {
	return instance.checkPermission(request, func(a access.Access) bool {
		return a.HasReadPermission()
	})
}

func (instance *RPC) hasWritePermission(request *restful.Request) bool {
	return instance.checkPermission(request, func(a access.Access) bool {
		return a.HasWritePermission()
	})
}

func (instance *RPC) onReadPermission(request *restful.Request, response *restful.Response, doThis func()) {
	if instance.hasReadPermission(request) {
		doThis()
	} else {
		response.WriteErrorString(http.StatusForbidden, "No read permission to instance endpoint.")
	}
}

func (instance *RPC) onWritePermission(request *restful.Request, response *restful.Response, doThis func()) {
	if instance.hasWritePermission(request) {
		doThis()
	} else {
		response.WriteErrorString(http.StatusForbidden, "No write permission to instance endpoint.")
	}
}

func (instance *RPC) config(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		response.WriteEntity(instance.caretakerd.ConfigObject())
	})
}

func (instance *RPC) controlConfig(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		response.WriteEntity(instance.caretakerd.Control().ConfigObject())
	})
}

func (instance *RPC) services(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		information := instance.execution.Information()
		response.WriteEntity(information)
	})
}

func (instance *RPC) service(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		serviceName := request.PathParameter("serviceName")
		services := instance.caretakerd.Services()
		if service := services.Get(serviceName); service != nil {
			information := instance.execution.InformationFor(service)
			response.WriteEntity(information)
		} else {
			response.WriteError(http.StatusNotFound, errors.New("Service '%s' does not exist.", serviceName))
		}
	})
}

func (instance *RPC) serviceConfig(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		instance.doWithService(request, response, func(service *service.Service) {
			response.WriteEntity(service.Config())
		})
	})
}

func (instance *RPC) serviceStatus(request *restful.Request, response *restful.Response) {
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

func (instance *RPC) servicePid(request *restful.Request, response *restful.Response) {
	instance.onReadPermission(request, response, func() {
		instance.doWithExecution(request, response, func(execution *service.Execution) {
			if execution != nil {
				response.Write([]byte(strconv.Itoa(execution.PID())))
			} else {
				response.Write([]byte("0"))
			}
		})
	})
}

func (instance *RPC) serviceRestart(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(service *service.Service) {
			err := instance.execution.Restart(service)
			if err == nil {
				response.Write([]byte("OK"))
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: "+err.Error())
			}
		})
	})
}

func (instance *RPC) serviceStart(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			err := instance.execution.Start(sc)
			if err == nil {
				response.Write([]byte("OK"))
			} else if sde, ok := err.(service.AlreadyRunningError); ok {
				response.WriteErrorString(http.StatusConflict, "ERROR: "+sde.Error())
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: "+err.Error())
			}
		})
	})
}

func (instance *RPC) serviceStop(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			err := instance.execution.Stop(sc)
			if err == nil {
				response.Write([]byte("OK"))
			} else if sde, ok := err.(service.AlreadyStoppedError); ok {
				response.WriteErrorString(http.StatusConflict, "ERROR: "+sde.Error())
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: "+err.Error())
			}
		})
	})
}

func (instance *RPC) serviceKill(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			err := instance.execution.Kill(sc)
			if err == nil {
				response.Write([]byte("OK"))
			} else if sde, ok := err.(service.AlreadyStoppedError); ok {
				response.WriteErrorString(http.StatusConflict, "ERROR: "+sde.Error())
			} else {
				response.WriteErrorString(http.StatusInternalServerError, "ERROR: "+err.Error())
			}
		})
	})
}

// SignalBody is a response structure.
type SignalBody struct {
	Signal values.Signal `json:"signal"`
}

func (instance *RPC) serviceSignal(request *restful.Request, response *restful.Response) {
	instance.onWritePermission(request, response, func() {
		instance.doWithService(request, response, func(sc *service.Service) {
			sb := SignalBody{}
			err := request.ReadEntity(&sb)
			if err != nil {
				response.WriteErrorString(http.StatusBadRequest, "ERROR: Illegal body. "+err.Error())
			} else {
				err = instance.execution.Kill(sc)
				if err == nil {
					response.Write([]byte("OK"))
				} else if sde, ok := err.(service.AlreadyStoppedError); ok {
					response.WriteErrorString(http.StatusConflict, "ERROR: "+sde.Error())
				} else {
					response.WriteErrorString(http.StatusInternalServerError, "ERROR: "+err.Error())
				}
			}
		})
	})
}

func (instance *RPC) doWithService(request *restful.Request, response *restful.Response, what func(service *service.Service)) {
	serviceName := request.PathParameter("serviceName")
	services := instance.caretakerd.Services()
	if service := services.Get(serviceName); service != nil {
		what(service)
	} else {
		response.WriteError(http.StatusNotFound, errors.New("Service '%s' does not exist.", serviceName))
	}
}

func (instance *RPC) doWithExecution(request *restful.Request, response *restful.Response, what func(execution *service.Execution)) {
	instance.doWithService(request, response, func(s *service.Service) {
		if e, ok := instance.execution.GetFor(s); ok {
			what(e)
		} else {
			what(nil)
		}
	})
}
