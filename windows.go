//go:build windows
// +build windows

package u

import (
	"golang.org/x/sys/windows"
)

func DisableCoreDump() error {
	windows.SetErrorMode(windows.SEM_NOGPFAULTERRORBOX)
	return nil
}
