package lib

import (
	"runtime"
)

// RootUser returns the name of the "root" / Admin user
func RootUser() string {
	if runtime.GOOS == "windows" {
		return "Administrator"
	}

	return "root"
}
