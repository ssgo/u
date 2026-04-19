//go:build linux
// +build linux

package u

import (
	"golang.org/x/sys/unix"
)

func DisableCoreDump() error {
	return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
}
