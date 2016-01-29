package values

import (
	"fmt"
	"github.com/echocat/caretakerd/errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var uriPattern = regexp.MustCompile("^([a-zA-Z0-9]+)://(.*)$")

// # Description
//
// This represents a socket address in format ``<protocol>://<target>``.
//
// # Protocols
//
// * ## ``tcp``
// This address connects or binds to a TCP socket. The ``target`` should be of format ``<host>:<port>``.
//
// * ## ``unix``
// This address connects or binds to a UNIX file socket. The ``target`` should be the location of the socket file.
type SocketAddress struct {
	Protocol Protocol
	Target   string
	Port     int
}

func (instance SocketAddress) AsScheme() string {
	return instance.Protocol.String()
}

func (instance SocketAddress) AsAddress() string {
	s, err := instance.checkedStringWithoutProtocol()
	if err != nil {
		panic(err)
	}
	return s
}

func (instance SocketAddress) String() string {
	s, err := instance.CheckedString()
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

func (instance SocketAddress) CheckedString() (string, error) {
	s, err := instance.checkedStringWithoutProtocol()
	if err != nil {
		return "", err
	}
	return instance.Protocol.String() + "://" + s, nil
}

func (instance SocketAddress) checkedStringWithoutProtocol() (string, error) {
	switch instance.Protocol {
	case Tcp:
		if !isValidPort(instance.Port) {
			return "", errors.New("Illegal port for protocol %v: %d", instance.Protocol, instance.Port)
		}
		if err := validateHost(instance.Target); err != nil {
			return "", errors.New("Illegal host for protocol %v: %s", instance.Protocol, instance.Target)
		}
		return fmt.Sprintf("%s:%d", instance.Target, instance.Port), nil
	case Unix:
		if instance.Port != 0 {
			return "", errors.New("For protocol %v is no port allowed.", instance.Protocol)
		}
		if len(strings.TrimSpace(instance.Target)) == 0 {
			return "", errors.New("For protocol %v is no target file defined.", instance.Protocol)
		}
		return fmt.Sprintf("%s", instance.Target), nil
	}
	return "", errors.New("Unknown protocol: %v", instance.Protocol)
}

func (instance *SocketAddress) Set(value string) error {
	match := uriPattern.FindStringSubmatch(value)
	if match != nil && len(match) == 3 {
		var protocol Protocol
		err := protocol.Set(match[1])
		if err != nil {
			return err
		}
		switch protocol {
		case Tcp:
			return instance.SetTcp(match[2])
		case Unix:
			return instance.SetUnix(match[2])
		}
		return errors.New("Unknown protocol %v in address '%v'.", protocol, value)
	} else {
		return errors.New("Illegal socket address: %s", value)
	}
}

func (instance *SocketAddress) SetTcp(value string) error {
	lastDoubleDot := strings.LastIndex(value, ":")
	if lastDoubleDot <= 0 || lastDoubleDot+2 >= len(value) {
		return errors.New("No port specified for address '%v'.", value)
	}
	host := value[:lastDoubleDot]
	plainPort := value[lastDoubleDot+1:]
	port, err := strconv.Atoi(plainPort)
	if err != nil || !isValidPort(port) {
		return errors.New("'%v' of specified address '%v' is not a valid port number.", plainPort, value)
	}
	if err := validateHost(host); err != nil {
		return errors.New("'%v' of specified address '%v' is not a valid host.", host, value)
	}
	(*instance).Protocol = Tcp
	(*instance).Target = host
	(*instance).Port = port
	return nil
}

func (instance *SocketAddress) SetUnix(value string) error {
	(*instance).Protocol = Unix
	(*instance).Target = value
	(*instance).Port = 0
	return nil
}

func (instance SocketAddress) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

func (instance *SocketAddress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance SocketAddress) Validate() error {
	_, err := instance.CheckedString()
	return err
}
