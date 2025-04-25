// Package gosafetmp implements temporal directories creation with
// solid guaranties that they will be deleted after program finish
// one way or another.
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

// Debug flags
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

// Destroy removes path and all nested entities (if there are).
// Does NOT return error if path dont exists.
// Currently just calls os.RemoveAll(path).
// But may implement more sophisticated strategies (like shred) in future.
func Destroy(path string) error {
	return os.RemoveAll(path)
}

// TmpDirManager is a factory that creates tmp dirs themselves.
// It implements multiple strategies to make sure that
// all children dirs will be deleted.
type TmpDirManager struct {
	baseDir string
}

// Cleanup removes all tmp dirs, created with TmpDirManager
func (t TmpDirManager) Cleanup() error {
	return Destroy(t.baseDir)
}

// GetBaseDir returns base directory where all children tmp dirs are created.
func (t TmpDirManager) GetBaseDir() string {
	return t.baseDir
}

// IsInTMPFS reports whether TmpDirManager's basic directory is inside
// inmemory filesystem. On MS Windows returns false all the time.
func (t TmpDirManager) IsInTMPFS() bool {
	return IsInTMPFS(t.baseDir)
}

// NewDir generates new tmp dir and returns path to it.
// It is safe to call NewDir concurrently.
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
		markForAutoDelete(basedir)
	}
	lockfile := path.Join(basedir, "lock")
	lockFile(lockfile)
	dirman := TmpDirManager{
		baseDir: basedir,
	}
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

// Setup spawn instance of TmpDirManager.
// This function returns same instance each time called.
func Setup() (*TmpDirManager, error) {
	once.Do(func() {
		instance, setupErr = setupOnce()
	})
	return instance, setupErr
}
