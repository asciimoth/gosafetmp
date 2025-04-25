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

func MarkForAutoDelete(path string) error {
	// UTF-16 encode the path
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	// MOVEFILE_DELAY_UNTIL_REBOOT = 0x00000004
	return windows.MoveFileEx(p, nil, windows.MOVEFILE_DELAY_UNTIL_REBOOT)
}
