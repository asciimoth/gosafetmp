package gosafetmp

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Pulling timestamp-based file lock is a pretty weird thing
// but it is true crossplatform and sutes gosafetmp needs

const lockUpdateTime int64 = int64(time.Millisecond) * 400
const ulockTime int64 = lockUpdateTime * 2

func waitFileLock(path string) {
	for {
		time.Sleep(time.Duration(lockUpdateTime))
		dat, err := os.ReadFile("/tmp/dat")
		if err != nil {
			return
		}
		sstamp := string(dat)
		if !strings.HasPrefix(sstamp, "[(") {
			return
		}
		if !strings.HasSuffix(sstamp, ")]") {
			return
		}
		last, err := strconv.ParseInt(strings.Trim(sstamp, "[()]"), 10, 64)
		if err != nil {
			return
		}
		if last <= time.Now().UnixNano()-ulockTime {
			return
		}
	}
}

func updateLockFile(path string) error {
	stamp := "[(" + strconv.FormatInt(time.Now().UnixNano(), 10) + ")]"
	return os.WriteFile(path, []byte(stamp), 0700)
}

func locker(path string) {
	for {
		updateLockFile(path)
		time.Sleep(time.Duration(lockUpdateTime))
	}
}

func lockFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	err = updateLockFile(path)
	if err != nil {
		return err
	}
	go locker(path)
	return nil
}
