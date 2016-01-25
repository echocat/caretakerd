package values

import (
    "strings"
    "github.com/echocat/caretakerd/errors"
    "regexp"
    "fmt"
    "strconv"
    "net"
)

var uriPattern = regexp.MustCompile("^([a-zA-Z0-9]+)://(.*)$")

// @id SocketAddress
// @type simple
//
// ## Description
//
// This represents a socket address in format ``<protocol>://<target>``.
//
// ## Protocols
//
// ### ``tcp``
//
// This address connects or binds to a TCP socket. The ``target`` should be of format ``<host>:<port>``.
//
// ### ``unix``
//
// This address connects or binds to a UNIX file socket. The ``target`` should be the location of the socket file.
type SocketAddress struct {
    Protocol Protocol
    Target   string
    Port     int
}

func (this SocketAddress) AsScheme() string {
    return this.Protocol.String()
}

func (this SocketAddress) AsAddress() string {
    s, err := this.checkedStringWithoutProtocol()
    if err != nil {
        panic(err)
    }
    return s
}

func (this SocketAddress) String() string {
    s, err := this.CheckedString()
    if err != nil {
        panic(err)
    }
    return s
}

func isValidPort(number int) bool {
    return number > 0 && number < 65535
}

func validateHost(host string) error {
    _, err := net.ResolveIPAddr("", host)
    return err
}

func (this SocketAddress) CheckedString() (string, error) {
    s, err := this.checkedStringWithoutProtocol()
    if err != nil {
        return "", err
    }
    return this.Protocol.String() + "://" + s, nil
}

func (this SocketAddress) checkedStringWithoutProtocol() (string, error) {
    switch this.Protocol {
    case Tcp:
        if !isValidPort(this.Port) {
            return "", errors.New("Illegal port for protocol %v: %d", this.Protocol, this.Port)
        }
        if err := validateHost(this.Target); err != nil {
            return "", errors.New("Illegal host for protocol %v: %s", this.Protocol, this.Target)
        }
        return fmt.Sprintf("%s:%d", this.Target, this.Port), nil
    case Unix:
        if this.Port != 0 {
            return "", errors.New("For protocol %v is no port allowed.", this.Protocol)
        }
        if len(strings.TrimSpace(this.Target)) == 0 {
            return "", errors.New("For protocol %v is no target file defined.", this.Protocol)
        }
        return fmt.Sprintf("%s", this.Target), nil
    }
    return "", errors.New("Unknown protocol: %v", this.Protocol)
}

func (this *SocketAddress) Set(value string) error {
    match := uriPattern.FindStringSubmatch(value)
    if match != nil && len(match) == 3 {
        var protocol Protocol
        err := protocol.Set(match[1])
        if err != nil {
            return err
        }
        switch protocol {
        case Tcp:
            return this.SetTcp(match[2])
        case Unix:
            return this.SetUnix(match[2])
        }
        return errors.New("Unknown protocol %v in address '%v'.", protocol, value)
    } else {
        return errors.New("Illegal socket address: %s", value)
    }
}

func (this *SocketAddress) SetTcp(value string) error {
    lastDoubleDot := strings.LastIndex(value, ":")
    if lastDoubleDot <= 0 || lastDoubleDot + 2 >= len(value) {
        return errors.New("No port specified for address '%v'.", value)
    }
    host := value[:lastDoubleDot]
    plainPort := value[lastDoubleDot + 1:]
    port, err := strconv.Atoi(plainPort)
    if err != nil || !isValidPort(port) {
        return errors.New("'%v' of specified address '%v' is not a valid port number.", plainPort, value)
    }
    if err := validateHost(host); err != nil {
        return errors.New("'%v' of specified address '%v' is not a valid host.", host, value)
    }
    (*this).Protocol = Tcp
    (*this).Target = host
    (*this).Port = port
    return nil
}

func (this *SocketAddress) SetUnix(value string) error {
    (*this).Protocol = Unix
    (*this).Target = value
    (*this).Port = 0
    return nil
}

func (this SocketAddress) MarshalYAML() (interface{}, error) {
    return this.String(), nil
}

func (this *SocketAddress) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this SocketAddress) Validate() error {
    _, err := this.CheckedString()
    return err
}
