package client

import (
    "gopkg.in/jmcvetta/napping.v3"
    "net/http"
    "crypto/tls"
    "crypto/x509"
    "net"
    "time"
    "io/ioutil"
    "strings"
    . "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/errors"
    "github.com/echocat/caretakerd/config"
    "github.com/echocat/caretakerd/service"
    sconfig "github.com/echocat/caretakerd/service/config"
    "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/service/signal"
    "github.com/echocat/caretakerd/control"
    "github.com/echocat/caretakerd/service/status"
)

type AccessDeniedError struct {
    url string
}

func (this AccessDeniedError) Error() string {
    return "Access to " + this.url + " is denied."
}

type ConflictError struct {
    error string
}

func (this ConflictError) Error() string {
    return this.error
}

type ServiceNotFoundError struct {}

func (this ServiceNotFoundError) Error() string {
    return "Service not found."
}

type ClientFactory struct {
    config *config.Config
}

func NewClientFactory(config *config.Config) *ClientFactory {
    return &ClientFactory{
        config: config,
    }
}

func (this *ClientFactory) NewClient() (*Client, error) {
    return NewClient(this.config)
}

type Client struct {
    address SocketAddress
    session *napping.Session
}

func NewClient(config *config.Config) (*Client, error) {
    session, err := sessionFor(config)
    if err != nil {
        return nil, err
    }
    return &Client{
        address: config.Rpc.Listen,
        session: session,
    }, nil
}

func sessionFor(config *config.Config) (*napping.Session, error) {
    httpClient, err := httpClientFor(config)
    if err != nil {
        return nil, err
    }
    return &napping.Session{
        Client: httpClient,
    }, nil
}

func httpClientFor(config *config.Config) (*http.Client, error) {
    transport, err := transportFor(config)
    if err != nil {
        return nil, err
    }
    return &http.Client{
        Transport: transport,
    }, nil
}

func transportFor(config *config.Config) (*http.Transport, error) {
    tlsConfig, err := tlsConfigFor(config)
    if err != nil {
        return nil, err
    }
    return &http.Transport{
        DialTLS: func(network, addr string) (net.Conn, error) {
            return dialTlsWithOwnChecks(config, tlsConfig)
        },
        TLSClientConfig: tlsConfig,
    }, nil
}

func tlsConfigFor(config *config.Config) (*tls.Config, error) {
    certificates, err := parseCertificatesInFile(config.Control.Access.PemFile)
    if err != nil {
        return nil, err
    }
    certificatePool, err := certPoolFor(certificates)
    if err != nil {
        return nil, err
    }
    return &tls.Config{
        Certificates: certificates,
        InsecureSkipVerify: true,
        RootCAs: certificatePool,
    }, nil
}

func parseCertificatesInFile(filename String) ([]tls.Certificate, error) {
    fileContent, err := ioutil.ReadFile(filename.String())
    if err != nil {
        return nil, errors.New("Could not read pem file '%v'.", filename).CausedBy(err)
    }
    return parseCertificates(fileContent)
}

func parseCertificates(pem []byte) ([]tls.Certificate, error) {
    cert, err := tls.X509KeyPair(pem, pem)
    if err != nil {
        return []tls.Certificate{}, err
    }
    return []tls.Certificate{cert}, nil
}

func certPoolFor(certificates []tls.Certificate) (*x509.CertPool, error) {
    result := x509.NewCertPool()
    for _, certificate := range certificates {
        for _, plainCert := range certificate.Certificate {
            cert, err := x509.ParseCertificate(plainCert)
            if err != nil {
                return nil, errors.New("Cannot load certificate from keypair.").CausedBy(err)
            }
            result.AddCert(cert)
        }
    }
    return result, nil
}

func dialTlsWithOwnChecks(config *config.Config, tlsConfig *tls.Config) (net.Conn, error) {
    var err error
    var tlsConn *tls.Conn

    address := config.Rpc.Listen
    tlsConn, err = tls.Dial(address.AsScheme(), address.AsAddress(), tlsConfig)
    if err != nil {
        return nil, err
    }

    if err = tlsConn.Handshake(); err != nil {
        tlsConn.Close()
        return nil, err
    }

    opts := x509.VerifyOptions{
        Roots:         tlsConfig.RootCAs,
        CurrentTime:   time.Now(),
        DNSName:       "",
        Intermediates: x509.NewCertPool(),
    }

    certs := tlsConn.ConnectionState().PeerCertificates
    for i, cert := range certs {
        if i == 0 {
            continue
        }
        opts.Intermediates.AddCert(cert)
    }

    _, err = certs[0].Verify(opts)
    if err != nil {
        tlsConn.Close()
        return nil, err
    }

    return tlsConn, err
}

func (this *Client) GetConfig() (config.Config, error) {
    target := config.Config{}
    err := this.get("config", &target)
    if err != nil {
        return config.Config{}, err
    }
    return target, nil
}

func (this *Client) GetControlConfig() (control.Config, error) {
    target := control.Config{}
    err := this.get("control/config", &target)
    if err != nil {
        return control.Config{}, err
    }
    return target, nil
}

func (this *Client) GetServices() (map[string]service.Information, error) {
    target := map[string]service.Information{}
    err := this.get("services", &target)
    if err != nil {
        return map[string]service.Information{}, err
    }
    return target, nil
}

func (this *Client) GetService(name string) (service.Information, error) {
    target := service.Information{}
    err := this.get("service/" + name, &target)
    if err != nil {
        return service.Information{}, err
    }
    return target, nil
}

func (this *Client) GetServiceConfig(name string) (sconfig.Config, error) {
    target := sconfig.Config{}
    err := this.get("service/" + name + "/config", &target)
    if err != nil {
        return sconfig.Config{}, err
    }
    return target, nil
}

func (this *Client) GetServiceStatus(name string) (status.Status, error) {
    var target status.Status
    plainTarget, err := this.getPlain("service/" + name + "/status")
    if err != nil {
        return target, err
    }
    err = target.Set(plainTarget)
    if err != nil {
        return target, err
    }
    return target, nil
}

func (this *Client) GetServicePid(name string) (values.Integer, error) {
    var target values.Integer
    plainTarget, err := this.getPlain("service/" + name + "/pid")
    if err != nil {
        return target, err
    }
    err = target.Set(plainTarget)
    if err != nil {
        return target, err
    }
    return target, nil
}

func (this *Client) StartService(name string) (error) {
    err := this.post("service/" + name + "/start", nil)
    if _, ok := err.(ConflictError); ok {
       return ConflictError{error: "Service '" + name + "' is already running."}
    }
    return err
}

func (this *Client) RestartService(name string) (error) {
    return this.post("service/" + name + "/restart", nil)
}

func (this *Client) StopService(name string) (error) {
    err := this.post("service/" + name + "/stop", nil)
    if _, ok := err.(ConflictError); ok {
        return ConflictError{error: "Service '" + name + "' is down."}
    }
    return err
}

func (this *Client) KillService(name string) (error) {
    err := this.post("service/" + name + "/kill", nil)
    if _, ok := err.(ConflictError); ok {
        return ConflictError{error: "Service '" + name + "' is down."}
    }
    return err
}

func (this *Client) SignalService(name string, s signal.Signal) (error) {
    payload := map[string]string{
        "signal": s.String(),
    }
    err := this.post("service/" + name + "/signal", &payload)
    if _, ok := err.(ConflictError); ok {
        return ConflictError{error: "Service '" + name + "' is down."}
    }
    return err
}

func (this *Client) get(path string, target interface{}) error {
    resp, err := this.session.Get("https://caretakerd/" + path, nil, target, nil)
    if err != nil {
        return err
    }
    if resp.Status() != 200 {
        return errors.New("Unexpected response from remote %v: %d - %s", this.address, resp.Status(), resp.RawText())
    }
    return nil
}

func (this *Client) transformError(path string, resp *napping.Response, err error) error {
    if err != nil {
        return err
    }
    if resp.Status() == http.StatusForbidden {
        return AccessDeniedError{url: resp.Url}
    }
    if resp.Status() == http.StatusConflict {
        return ConflictError{error: path + " is conflict state."}
    }
    if resp.Status() == http.StatusNotFound {
        body := resp.RawText()
        if strings.HasPrefix(body, "Service '") && strings.HasSuffix(body, "' does not exist.") {
            return ServiceNotFoundError{}
        }
    }
    if resp.Status() != http.StatusOK {
        return errors.New("Unexpected response from '%v': %d - %s", this.address, resp.Status(), resp.RawText())
    }
    return nil
}

func (this *Client) getPlain(path string) (string, error) {
    resp, err := this.session.Get("https://caretakerd/" + path, nil, nil, nil)
    targetErr := this.transformError(path, resp, err)
    if targetErr != nil {
        return "", targetErr
    }
    return resp.RawText(), nil
}

func (this *Client) post(path string, payload interface{}) error {
    resp, err := this.session.Post("https://caretakerd/" + path, payload, nil, nil)
    return this.transformError(path, resp, err)
}
