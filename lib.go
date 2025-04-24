package gosafetmp

import (
	"os"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	SHOULD_SPAWN_REAPER = true
)

// There are MUST be only one instance of TmpDirManager in whole program
var (
	instance *TmpDirManager = nil
	setupErr error          = nil
	once     sync.Once
	counter  atomic.Int64
)

func Destroy(path string) error {
	return os.RemoveAll(path)
}

type TmpDirManager struct {
	baseDir string
}

func (t TmpDirManager) Cleanup() error {
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
	if SHOULD_SPAWN_REAPER {
		checkReaper()
	}
	basedir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	// TODO: Setup basedir autodeletion (in Windows case)
	lockfile := path.Join(basedir, "lock")
	lockFile(lockfile)
	dirman := TmpDirManager{
		baseDir: basedir,
	}
	// TODO: Set up OS signals handlers
	if SHOULD_SPAWN_REAPER {
		err := spawnReaper(basedir, lockfile)
		if err != nil {
			return nil, err
		}
	}
	return &dirman, nil
}

func Setup() (*TmpDirManager, error) {
	once.Do(func() {
		instance, setupErr = setupOnce()
	})
	return instance, setupErr
}
