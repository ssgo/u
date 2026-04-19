//go:build darwin
// +build darwin

package u

import (
	"golang.org/x/sys/unix"
)

func DisableCoreDump() error {
	var rlimit unix.Rlimit
	rlimit.Cur = 0
	rlimit.Max = 0
	return unix.Setrlimit(unix.RLIMIT_CORE, &rlimit)
}
