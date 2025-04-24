package gosafetmp

import (
	"os"
	"path"
	"strconv"
	"sync"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

// There are MUST be only one instance of TmpDirManager in whole program
var (
	instance *TmpDirManager = nil
	once     sync.Once
	counter  atomic.Int64
)

const TMPFS_MAGIC = 0x01021994

func IsInTMPFS(path string) bool {
	// TODO: On non UNIX systems just return false
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return false
	}
	return st.Type == TMPFS_MAGIC
}

func Destroy(path string) error {
	// TODO: Maybe also use GNU shred on linux systems?
	return os.RemoveAll(path)
}

type TmpDirManager struct {
	baseDir string
}

func (t TmpDirManager) Cleanup() error {
	// TODO: Release lock file
	return Destroy(t.baseDir)
}

func (t TmpDirManager) GetBaseDir() string {
	return t.baseDir
}

func (t TmpDirManager) IsInTMPFS() bool {
	return IsInTMPFS(t.baseDir)
}

func (t TmpDirManager) NewDir() (string, error) {
	id := counter.Add(1)
	dir := path.Join(t.baseDir, strconv.FormatInt(id, 16))
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", nil
	}
	return dir, nil
}

func setupOnce() (*TmpDirManager, error) {
	// TODO: Check if it is reaper daemon
	basedir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	// TODO: Setup basedir autodeletion (in Windows case)
	// TODO: Setup lock file in basedir
	dirman := TmpDirManager{
		baseDir: basedir,
	}
	// TODO: Set up OS signals handlers
	// TODO: Spawn reaper daemon
	return &dirman, nil
}

func Setup() (*TmpDirManager, error) {
	var err error = nil
	once.Do(func() {
		instance, err = setupOnce()
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
