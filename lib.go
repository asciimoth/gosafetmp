package gosafetmp

import (
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	SHOULD_SPAWN_REAPER         = true
	SHOULD_MARK_FOR_AUTO_DELETE = true
	SHOULD_CATCH_SIGNALS        = true
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

func catchSignals(callback func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		os.Interrupt,    // CTRL+C on both UNIX & Windows
		syscall.SIGTERM, // “kill” on UNIX, service stop on Windows
		syscall.SIGHUP,  // hangup on UNIX
	)
	go func() {
		<-sigs
		callback()
		os.Exit(1)
	}()
}

func setupOnce() (*TmpDirManager, error) {
	if SHOULD_SPAWN_REAPER {
		checkReaper()
	}
	basedir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	if SHOULD_MARK_FOR_AUTO_DELETE {
		MarkForAutoDelete(basedir)
	}
	lockfile := path.Join(basedir, "lock")
	lockFile(lockfile)
	dirman := TmpDirManager{
		baseDir: basedir,
	}
	// TODO: Set up OS signals handlers
	if SHOULD_CATCH_SIGNALS {
		catchSignals(func() { Destroy(basedir) })
	}
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
