package defaults

import (
	"github.com/echocat/caretakerd/values"
)

func ListenAddress() values.SocketAddress {
	return values.SocketAddress{
		Protocol: values.Tcp,
		Target: "localhost",
		Port: 57955,
	}
}

