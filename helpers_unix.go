//go:build !windows

package gosafetmp

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const tmpfs_magic = 0x01021994

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

// IsInTMPFS reports whether path is inside inmemory filesystem.
// On MS Windows returns false all the time.
func IsInTMPFS(path string) bool {
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return false
	}
	return st.Type == tmpfs_magic
}

func markForAutoDelete(path string) error {
	// Do nothing
	// It is Window-only feature
	return nil
}
