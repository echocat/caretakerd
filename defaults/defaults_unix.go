// +build linux,darwin

package defaults

import (
	. "github.com/echocat/caretakerd/values"
)

func AuthFileKeyFilename() String {
	return String("/var/run/caretakerd.key")
}

func ConfigFilename() String {
	return String("/etc/caretakerd.yaml")
}
