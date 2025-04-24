package gosafetmp

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const basenv = "__GOSAFETMP_BASE__"
const lockenv = "__GOSAFETMP_LOCK__"

func checkReaper() {
	basedir := ""
	lockfile := ""
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == basenv {
			basedir = pair[1]
		} else if pair[0] == lockenv {
			lockfile = pair[1]
		}
	}
	if basedir == "" || lockfile == "" {
		return
	}
	waitFileLock(lockfile)
	Destroy(basedir)
	os.Exit(0)
}

func spawnReaper(basedir, lockfile string) error {
	cmd := exec.Command(os.Args[0])

	cmd.Env = append(os.Environ(),
		basenv+"="+basedir,
		lockenv+"="+lockfile,
	)

	cmd.SysProcAttr = sysProcAttr()

	if runtime.GOOS != "windows" {
		// Direct stdio to /dev/null
		null, _ := os.OpenFile("/dev/null", os.O_RDWR, 0)
		cmd.Stdin = null
		cmd.Stdout = null
		cmd.Stderr = null
	}

	return cmd.Start()
}
