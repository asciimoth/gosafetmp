//go:build !windows

package gosafetmp

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const TMPFS_MAGIC = 0x01021994

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

func IsInTMPFS(path string) bool {
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return false
	}
	return st.Type == TMPFS_MAGIC
}

func MarkForAutoDelete(path string) error {
	// Do nothing
	// It is Window-only feature
	return nil
}
