package defaults

import (
	"github.com/echocat/caretakerd/values"
	"os"
	"runtime"
)

type Defaults struct {
	ListenAddress       values.SocketAddress
	AuthFileKeyFilename values.String
	ConfigFilename      values.String
}

var listenAddress = values.SocketAddress{
	Protocol: values.TCP,
	Target:   "localhost",
	Port:     57955,
}

const unixAuthFileKeyFilename = values.String("/var/run/caretakerd.key")
const windowsAuthFileKeyFilename = values.String("C:\\ProgramData\\caretakerd\\access.key")
const unixConfigFilename = values.String("/etc/caretakerd.yaml")
const windowsConfigFilename = values.String("C:\\ProgramData\\caretakerd\\config.yaml")

var allDefaults = map[string]Defaults{
	"linux": Defaults{
		ListenAddress:       listenAddress,
		AuthFileKeyFilename: unixAuthFileKeyFilename,
		ConfigFilename:      unixConfigFilename,
	},
	"windows": Defaults{
		ListenAddress:       listenAddress,
		AuthFileKeyFilename: windowsAuthFileKeyFilename,
		ConfigFilename:      windowsConfigFilename,
	},
	"darwin": Defaults{
		ListenAddress:       listenAddress,
		AuthFileKeyFilename: unixAuthFileKeyFilename,
		ConfigFilename:      unixConfigFilename,
	},
}

func GetDefaults() Defaults {
	goos := os.Getenv("GOOS")
	if goos != "" {
		return GetDefaultsFor(goos)
	}
	return GetDefaultsFor(runtime.GOOS)
}

func GetDefaultsFor(name string) Defaults {
	if defaults, ok := allDefaults[name]; ok {
		return defaults
	}
	panic("Unsupported os: " + name)
}

func ListenAddress() values.SocketAddress {
	return GetDefaults().ListenAddress
}

func AuthFileKeyFilename() values.String {
	return GetDefaults().AuthFileKeyFilename
}

func ConfigFilename() values.String {
	return GetDefaults().ConfigFilename
}
