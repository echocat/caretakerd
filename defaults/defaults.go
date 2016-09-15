package defaults

import (
	"github.com/echocat/caretakerd/values"
	"os"
	"runtime"
)

// Defaults holds default values for a specific platform.
type Defaults struct {
	// ListenAddress represents the default address caretakerd listens to.
	ListenAddress values.SocketAddress
	// AuthFileKeyFilename represents the default file name caretakerd uses to store 
	// the key for the caretakerctl/control process.
	AuthFileKeyFilename values.String
	// ConfigFilename represents the default location where caretakerd searches for its config file (yaml).
	ConfigFilename values.String
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
	"linux": {
		ListenAddress:       listenAddress,
		AuthFileKeyFilename: unixAuthFileKeyFilename,
		ConfigFilename:      unixConfigFilename,
	},
	"windows": {
		ListenAddress:       listenAddress,
		AuthFileKeyFilename: windowsAuthFileKeyFilename,
		ConfigFilename:      windowsConfigFilename,
	},
	"darwin": {
		ListenAddress:       listenAddress,
		AuthFileKeyFilename: unixAuthFileKeyFilename,
		ConfigFilename:      unixConfigFilename,
	},
}

// GetDefaults queries the caretakerd default values of the current instance.
// This will be influenced by the GOOS environment variable.
func GetDefaults() Defaults {
	goos := os.Getenv("GOOS")
	if goos != "" {
		return GetDefaultsFor(goos)
	}
	return GetDefaultsFor(runtime.GOOS)
}

// GetDefaultsFor queries the caretakerd default values of the given platform.
func GetDefaultsFor(platform string) Defaults {
	if defaults, ok := allDefaults[platform]; ok {
		return defaults
	}
	panic("Unsupported os: " + platform)
}

// ListenAddress queries the ListenAddress of the current platform.
// This will be influenced by GOOS environment variable.
func ListenAddress() values.SocketAddress {
	return GetDefaults().ListenAddress
}

// AuthFileKeyFilename queries the AuthFileKeyFilename of the current platform.
// This will be influenced by GOOS environment variable.
func AuthFileKeyFilename() values.String {
	return GetDefaults().AuthFileKeyFilename
}

// ConfigFilename queries the ConfigFilename of the current platform.
// This will be influenced by GOOS environment variable.
func ConfigFilename() values.String {
	return GetDefaults().ConfigFilename
}
