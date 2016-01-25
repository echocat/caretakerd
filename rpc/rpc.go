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

func (this *Rpc) Start() {
    go this.Run()
}

func (this *Rpc) Run() {
    defer panics.DefaultPanicHandler()
    container := restful.NewContainer()

    ws := new(restful.WebService)
    ws.Produces(restful.MIME_JSON)

    ws.Route(ws.GET("/config").To(this.config))

    ws.Route(ws.GET("/control/config").To(this.controlConfig))

    ws.Route(ws.GET("/services").To(this.services))

    ws.Route(ws.GET("/service/{serviceName}").To(this.service))
    ws.Route(ws.GET("/service/{serviceName}/config").To(this.serviceConfig))
    ws.Route(ws.GET("/service/{serviceName}/state").To(this.serviceStatus))
    ws.Route(ws.GET("/service/{serviceName}/pid").To(this.servicePid))

    ws.Route(ws.POST("/service/{serviceName}/restart").To(this.serviceRestart))
    ws.Route(ws.POST("/service/{serviceName}/status").To(this.serviceStart))
    ws.Route(ws.POST("/service/{serviceName}/stop").To(this.serviceStop))
    ws.Route(ws.POST("/service/{serviceName}/kill").To(this.serviceKill))
    ws.Route(ws.POST("/service/{serviceName}/signal").To(this.serviceSignal))

    container.Add(ws)

    server := &http.Server{
        Handler: container,
        ErrorLog: log.New(this.logger.ReceiverFor(logger.Debug), "", 0),
    }
    this.logger.Log(logger.Debug, "Rpc will bind to %v...", this.conf.Listen)
    listener, err := net.Listen(this.conf.Listen.AsScheme(), this.conf.Listen.AsAddress())
    if err != nil {
        log.Fatal(err)
    }
    sl, err2 := NewStoppableListener(listener)
    if err2 != nil {
        panics.New("Could not create listener.").CausedBy(err2).Throw()
    }
    defer func() {
        (*this).listener = nil
    }()
    (*this).listener = sl
    if err := server.Serve(this.secure(sl)); err != nil {
        if _, ok := err.(ListenerStopped); !ok {
            panics.New("Could not listen.").CausedBy(err2).Throw()
        }
    }
}

func (this *Rpc) secure(in net.Listener) net.Listener {
    out := in
    sec := this.caretakerd.KeyStore()
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

func (this *Rpc) Stop() {
    listener := (*this).listener
    if listener != nil {
        listener.Close()
    }
}

func (this *Rpc) checkPermission(request *restful.Request, permissionChecker func(access.Access) bool) bool {
    if request != nil {
        hr := request.Request
        if hr != nil {
            cs := hr.TLS
            if cs != nil {
                for _, cert := range cs.PeerCertificates {
                    ctl := this.caretakerd.Control()
                    acc := ctl.Access()
                    if acc.IsCertValid(cert) && permissionChecker(*acc) {
                        return true
                    }
                    for _, serv := range (*this.caretakerd.Services()) {
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

func (this *Rpc) hasReadPermission(request *restful.Request) bool {
    return this.checkPermission(request, func(a access.Access) bool {
        return a.HasReadPermission()
    })
}

func (this *Rpc) hasWritePermission(request *restful.Request) bool {
    return this.checkPermission(request, func(a access.Access) bool {
        return a.HasWritePermission()
    })
}

func (this *Rpc) onReadPermission(request *restful.Request, response *restful.Response, doThis func()) {
    if this.hasReadPermission(request) {
        doThis()
    } else {
        response.WriteErrorString(http.StatusForbidden, "No read permission to this endpoint.")
    }
}

func (this *Rpc) onWritePermission(request *restful.Request, response *restful.Response, doThis func()) {
    if this.hasWritePermission(request) {
        doThis()
    } else {
        response.WriteErrorString(http.StatusForbidden, "No write permission to this endpoint.")
    }
}

func (this *Rpc) config(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        response.WriteEntity(this.caretakerd.ConfigObject())
    })
}

func (this *Rpc) controlConfig(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        response.WriteEntity(this.caretakerd.Control().ConfigObject())
    })
}

func (this *Rpc) services(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        information := this.execution.Information()
        response.WriteEntity(information)
    })
}

func (this *Rpc) service(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        serviceName := request.PathParameter("serviceName")
        services := this.caretakerd.Services()
        if service, ok := services.Get(serviceName); ok {
            information := this.execution.InformationFor(service)
            response.WriteEntity(information)
        } else {
            response.WriteError(http.StatusNotFound, errors.New("Service '%s' does not exist.", serviceName))
        }
    })
}

func (this *Rpc) serviceConfig(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        this.doWithService(request, response, func(service *service.Service) {
            response.WriteEntity(service.Config())
        })
    })
}

func (this *Rpc) serviceStatus(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        this.doWithExecution(request, response, func(execution *service.Execution) {
            if execution != nil {
                response.Write([]byte(execution.Status().String()))
            } else {
                response.Write([]byte(service.Down.String()))
            }
        })
    })
}

func (this *Rpc) servicePid(request *restful.Request, response *restful.Response) {
    this.onReadPermission(request, response, func() {
        this.doWithExecution(request, response, func(execution *service.Execution) {
            if execution != nil {
                response.Write([]byte(strconv.Itoa(execution.Pid())))
            } else {
                response.Write([]byte("0"))
            }
        })
    })
}

func (this *Rpc) serviceRestart(request *restful.Request, response *restful.Response) {
    this.onWritePermission(request, response, func() {
        this.doWithService(request, response, func(service *service.Service) {
            err := this.execution.Restart(service)
            if err == nil {
                response.Write([]byte("OK"))
            } else {
                response.WriteErrorString(http.StatusInternalServerError, "ERROR: " + err.Error())
            }
        })
    })
}

func (this *Rpc) serviceStart(request *restful.Request, response *restful.Response) {
    this.onWritePermission(request, response, func() {
        this.doWithService(request, response, func(sc *service.Service) {
            err := this.execution.Start(sc)
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

func (this *Rpc) serviceStop(request *restful.Request, response *restful.Response) {
    this.onWritePermission(request, response, func() {
        this.doWithService(request, response, func(sc *service.Service) {
            err := this.execution.Stop(sc)
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

func (this *Rpc) serviceKill(request *restful.Request, response *restful.Response) {
    this.onWritePermission(request, response, func() {
        this.doWithService(request, response, func(sc *service.Service) {
            err := this.execution.Kill(sc)
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

func (this *Rpc) serviceSignal(request *restful.Request, response *restful.Response) {
    this.onWritePermission(request, response, func() {
        this.doWithService(request, response, func(sc *service.Service) {
            sb := SignalBody{}
            err := request.ReadEntity(&sb)
            if err != nil {
                response.WriteErrorString(http.StatusBadRequest, "ERROR: Illegal body. " + err.Error())
            } else {
                err = this.execution.Kill(sc)
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

func (this *Rpc) doWithService(request *restful.Request, response *restful.Response, what func(service *service.Service)) {
    serviceName := request.PathParameter("serviceName")
    services := this.caretakerd.Services()
    if service, ok := services.Get(serviceName); ok {
        what(service)
    } else {
        response.WriteError(http.StatusNotFound, errors.New("Service '%s' does not exist.", serviceName))
    }
}

func (this *Rpc) doWithExecution(request *restful.Request, response *restful.Response, what func(execution *service.Execution)) {
    this.doWithService(request, response, func(s *service.Service) {
        if e, ok := this.execution.GetFor(s); ok {
            what(e)
        } else {
            what(nil)
        }
    })
}

