package defaults

import (
	"github.com/echocat/caretakerd/values"
	"os"
	"runtime"
)

// Defaults holds default values for a specific platform.
type Defaults struct {
	// ListenAddress is the address caretakerd will listen per default to.
	ListenAddress values.SocketAddress
	// AuthFileKeyFilename is the filename caretakerd will store by default
	// the key for caretakerctl/control process to it.
	AuthFileKeyFilename values.String
	// ConfigFilename is the default location where caretakerd searches for its config file (yaml).
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

// GetDefaults returns the caretakerd defaults for the current instance.
// This will be influenced by GOOS environment variable.
func GetDefaults() Defaults {
	goos := os.Getenv("GOOS")
	if goos != "" {
		return GetDefaultsFor(goos)
	}
	return GetDefaultsFor(runtime.GOOS)
}

// GetDefaultsFor returns the caretakerd defaults for given platform.
func GetDefaultsFor(platform string) Defaults {
	if defaults, ok := allDefaults[platform]; ok {
		return defaults
	}
	panic("Unsupported os: " + platform)
}

// ListenAddress returns the ListenAddress for the current platform.
// This will be influenced by GOOS environment variable.
func ListenAddress() values.SocketAddress {
	return GetDefaults().ListenAddress
}

// AuthFileKeyFilename returns the AuthFileKeyFilename for the current platform.
// This will be influenced by GOOS environment variable.
func AuthFileKeyFilename() values.String {
	return GetDefaults().AuthFileKeyFilename
}

// ConfigFilename returns the ConfigFilename for the current platform.
// This will be influenced by GOOS environment variable.
func ConfigFilename() values.String {
	return GetDefaults().ConfigFilename
}
