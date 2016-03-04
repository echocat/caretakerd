package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/control"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/values"
	"gopkg.in/jmcvetta/napping.v3"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// AccessDeniedError represents an error that occurs if someone tries to access a
// resource that he is not allowed to access.
type AccessDeniedError struct {
	url string
}

func (instance AccessDeniedError) Error() string {
	return "Access to " + instance.url + " is denied."
}

// ConflictError represents an error that occurs if someone tries to do an action
// on an entity that is in a different state.
type ConflictError struct {
	error string
}

func (instance ConflictError) Error() string {
	return instance.error
}

// ServiceNotFoundError represents an error that occurs if someone tries to access a
// service that does not exists.
type ServiceNotFoundError struct{}

func (instance ServiceNotFoundError) Error() string {
	return "Service not found."
}

// Factory is used to create new instances of the caretakerd Client.
type Factory struct {
	config *caretakerd.Config
}

// NewFactory creates a new instance of Factory.
func NewFactory(config *caretakerd.Config) *Factory {
	return &Factory{
		config: config,
	}
}

// NewClient creates a new Client.
func (instance *Factory) NewClient() (*Client, error) {
	return NewClient(instance.config)
}

// Client is used to access caretakerd remotely.
type Client struct {
	address values.SocketAddress
	session *napping.Session
}

// NewClient creates a new instance of Client with the given config.
func NewClient(config *caretakerd.Config) (*Client, error) {
	session, err := sessionFor(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		address: config.Rpc.Listen,
		session: session,
	}, nil
}

func sessionFor(config *caretakerd.Config) (*napping.Session, error) {
	httpClient, err := httpClientFor(config)
	if err != nil {
		return nil, err
	}
	return &napping.Session{
		Client: httpClient,
	}, nil
}

func httpClientFor(config *caretakerd.Config) (*http.Client, error) {
	transport, err := transportFor(config)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: transport,
	}, nil
}

func transportFor(config *caretakerd.Config) (*http.Transport, error) {
	tlsConfig, err := tlsConfigFor(config)
	if err != nil {
		return nil, err
	}
	return &http.Transport{
		DialTLS: func(network, addr string) (net.Conn, error) {
			return dialTLSWithOwnChecks(config, tlsConfig)
		},
		TLSClientConfig: tlsConfig,
	}, nil
}

func tlsConfigFor(config *caretakerd.Config) (*tls.Config, error) {
	certificates, err := parseCertificatesInFile(config.Control.Access.PemFile)
	if err != nil {
		return nil, err
	}
	certificatePool, err := certPoolFor(certificates)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       certificates,
		InsecureSkipVerify: true,
		RootCAs:            certificatePool,
	}, nil
}

func parseCertificatesInFile(filename values.String) ([]tls.Certificate, error) {
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

func dialTLSWithOwnChecks(config *caretakerd.Config, tlsConfig *tls.Config) (net.Conn, error) {
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

// GetConfig returns the config of the remote caretakerd instance.
func (instance *Client) GetConfig() (caretakerd.Config, error) {
	target := caretakerd.Config{}
	err := instance.get("config", &target)
	if err != nil {
		return caretakerd.Config{}, err
	}
	return target, nil
}

// GetControlConfig returns the control config of the remote caretakerd instance.
func (instance *Client) GetControlConfig() (control.Config, error) {
	target := control.Config{}
	err := instance.get("control/config", &target)
	if err != nil {
		return control.Config{}, err
	}
	return target, nil
}

// GetServices returns all services of the remote caretakerd instance.
func (instance *Client) GetServices() (map[string]service.Information, error) {
	target := map[string]service.Information{}
	err := instance.get("services", &target)
	if err != nil {
		return map[string]service.Information{}, err
	}
	return target, nil
}

// GetService returns the given service (by name) of the remote caretakerd instance.
func (instance *Client) GetService(name string) (service.Information, error) {
	target := service.Information{}
	err := instance.get("service/"+name, &target)
	if err != nil {
		return service.Information{}, err
	}
	return target, nil
}

// GetServiceConfig returns the given service config (by name) of the remote caretakerd instance.
func (instance *Client) GetServiceConfig(name string) (service.Config, error) {
	target := service.Config{}
	err := instance.get("service/"+name+"/config", &target)
	if err != nil {
		return service.Config{}, err
	}
	return target, nil
}

// GetServiceStatus returns the given service status (by name) of the remote caretakerd instance.
func (instance *Client) GetServiceStatus(name string) (service.Status, error) {
	var target service.Status
	plainTarget, err := instance.getPlain("service/" + name + "/status")
	if err != nil {
		return target, err
	}
	err = target.Set(plainTarget)
	if err != nil {
		return target, err
	}
	return target, nil
}

// GetServicePid returns the given service PID (by name) of the remote caretakerd instance.
func (instance *Client) GetServicePid(name string) (values.Integer, error) {
	var target values.Integer
	plainTarget, err := instance.getPlain("service/" + name + "/pid")
	if err != nil {
		return target, err
	}
	err = target.Set(plainTarget)
	if err != nil {
		return target, err
	}
	return target, nil
}

// StartService starts the given service (by name) of the remote caretakerd instance.
func (instance *Client) StartService(name string) error {
	err := instance.post("service/"+name+"/start", nil)
	if _, ok := err.(ConflictError); ok {
		return ConflictError{error: "Service '" + name + "' is already running."}
	}
	return err
}

// RestartService restarts the given service (by name) of the remote caretakerd instance.
func (instance *Client) RestartService(name string) error {
	return instance.post("service/"+name+"/restart", nil)
}

// StopService stops the given service (by name) of the remote caretakerd instance.
func (instance *Client) StopService(name string) error {
	err := instance.post("service/"+name+"/stop", nil)
	if _, ok := err.(ConflictError); ok {
		return ConflictError{error: "Service '" + name + "' is down."}
	}
	return err
}

// KillService kills the given service (by name) of the remote caretakerd instance.
func (instance *Client) KillService(name string) error {
	err := instance.post("service/"+name+"/kill", nil)
	if _, ok := err.(ConflictError); ok {
		return ConflictError{error: "Service '" + name + "' is down."}
	}
	return err
}

// SignalService sends the given signal to the given service (by name) of the remote caretakerd instance.
func (instance *Client) SignalService(name string, s values.Signal) error {
	payload := map[string]string{
		"signal": s.String(),
	}
	err := instance.post("service/"+name+"/signal", &payload)
	if _, ok := err.(ConflictError); ok {
		return ConflictError{error: "Service '" + name + "' is down."}
	}
	return err
}

func (instance *Client) get(path string, target interface{}) error {
	resp, err := instance.session.Get("https://caretakerd/"+path, nil, target, nil)
	if err != nil {
		return err
	}
	if resp.Status() != 200 {
		return errors.New("Unexpected response from remote %v: %d - %s", instance.address, resp.Status(), resp.RawText())
	}
	return nil
}

func (instance *Client) transformError(path string, resp *napping.Response, err error) error {
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
		return errors.New("Unexpected response from '%v': %d - %s", instance.address, resp.Status(), resp.RawText())
	}
	return nil
}

func (instance *Client) getPlain(path string) (string, error) {
	resp, err := instance.session.Get("https://caretakerd/"+path, nil, nil, nil)
	targetErr := instance.transformError(path, resp, err)
	if targetErr != nil {
		return "", targetErr
	}
	return resp.RawText(), nil
}

func (instance *Client) post(path string, payload interface{}) error {
	resp, err := instance.session.Post("https://caretakerd/"+path, payload, nil, nil)
	return instance.transformError(path, resp, err)
}
