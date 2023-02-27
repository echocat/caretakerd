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

// @serializedAs string
// # Description
//
// SocketAddress represents a socket address in the format “<protocol>://<target>“.
//
// # Protocols
//
//   - **“tcp“** This address connects or binds to a TCP socket. The “target“ should be of format “<host>:<port>“.<br>
//     Examples:
//   - “tcp://localhost:57955“: Listen on IPv4 and IPv6 local addresses
//   - “tcp://[::1]:57955“: Listen on IPv6 local address
//   - “tcp://0.0.0.0:57955“: Listen on all addresses - this includes IPv4 and IPv6
//   - “tcp://192.168.0.1:57955“: Listen on specific IPv4 address
//   - **“unix“** This address connects or binds to a UNIX file socket. The “target“ should be the location of the socket file.<br>
//     Example:
//   - “unix:///var/run/caretakerd.sock“
type SocketAddress struct {
	Protocol Protocol
	Target   string
	Port     int
}

// AsScheme returns a string that represents this SocketAddress as scheme.
func (instance SocketAddress) AsScheme() string {
	return instance.Protocol.String()
}

// AsAddress returns a string that represents this SocketAddress as socket address.
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

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance SocketAddress) CheckedString() (string, error) {
	s, err := instance.checkedStringWithoutProtocol()
	if err != nil {
		return "", err
	}
	return instance.Protocol.String() + "://" + s, nil
}

func (instance SocketAddress) checkedStringWithoutProtocol() (string, error) {
	switch instance.Protocol {
	case TCP:
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
		return instance.Target, nil
	}
	return "", errors.New("Unknown protocol: %v", instance.Protocol)
}

// Set sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *SocketAddress) Set(value string) error {
	match := uriPattern.FindStringSubmatch(value)
	if len(match) == 3 {
		var protocol Protocol
		err := protocol.Set(match[1])
		if err != nil {
			return err
		}
		switch protocol {
		case TCP:
			return instance.SetTCP(match[2])
		case Unix:
			return instance.SetUnix(match[2])
		}
		return errors.New("Unknown protocol %v in address '%v'.", protocol, value)
	}
	return errors.New("Illegal socket address: %s", value)
}

// SetTCP set this SocketAddress instance to the given TCP address (without leadning tcp:// scheme)
func (instance *SocketAddress) SetTCP(value string) error {
	lastDoubleDot := strings.LastIndex(value, ":")
	if lastDoubleDot <= 0 || lastDoubleDot+2 >= len(value) {
		return errors.New("No port specified for address '%v'.", value)
	}
	host := value[:lastDoubleDot]
	plainPort := value[lastDoubleDot+1:]
	port, err := strconv.Atoi(plainPort)
	if err != nil || !isValidPort(port) {
		return errors.New("'%v' of specified address '%v' is not a valid port number", plainPort, value)
	}
	if err := validateHost(host); err != nil {
		return errors.New("'%v' of specified address '%v' is not a valid host", host, value)
	}
	(*instance).Protocol = TCP
	(*instance).Target = host
	(*instance).Port = port
	return nil
}

// SetUnix set this SocketAddress instance to the given Unix socket file (without leadning unix:// scheme)
func (instance *SocketAddress) SetUnix(value string) error {
	(*instance).Protocol = Unix
	(*instance).Target = value
	(*instance).Port = 0
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance SocketAddress) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *SocketAddress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance SocketAddress) Validate() error {
	_, err := instance.CheckedString()
	return err
}
