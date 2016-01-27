// +build windows

package access

import "os"

const defaultFilePermission = FilePermission(os.FileMode(0600))
