// +build windows

package defaults

import (
	. "github.com/echocat/caretakerd/values"
)

func AuthFileKeyFilename() String {
	return String("C:\\ProgramData\\caretakerd\\access.key")
}

func ConfigFilename() String {
	return String("C:\\ProgramData\\caretakerd\\config.yaml")
}

