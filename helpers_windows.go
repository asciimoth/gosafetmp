//go:build windows

package gosafetmp

import (
	"syscall"

	"golang.org/x/sys/windows"
)

const TMPFS_MAGIC = 0x01021994

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP |
			windows.DETACHED_PROCESS,
	}
}

func IsInTMPFS(path string) bool {
	return false
}
