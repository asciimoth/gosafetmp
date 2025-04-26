[![Go Reference](https://pkg.go.dev/badge/github.com/asciimoth/gosafetmp.svg)](https://pkg.go.dev/github.com/asciimoth/gosafetmp)
# gosafetmp
Gosafetmp implements temporary directory creation with solid guarantees that they will be deleted after the program finishes, one way or another.

# usage
```go
package main

import "github.com/asciimoth/gosafetmp"

func main() {
	// gosafetmp.Setup() MUST be called as close to the program start as possible
	// Any code above it will be executed twice
	factory, err := gosafetmp.Setup()
	if err != nil {
		painc(err)
	}

	// Use this construction instead of just `defer factory.Cleanup()`
	// to handle panics
	defer func() { factory.Cleanup() }()

	doWork(factory)
}

func doWork(fact gosafetmp.TmpDirManager) {
	tmpDir, err := fact.NewDir()

	// All dirs created with the factory will be deleted when the program finishes anyway
	// However, it's good practice to destroy them immediately after they're no longer needed
	defer gosafetmp.Destroy(tmpDir)

	/* ... Do some work that requires tmp files here ... */
}
```

# strategies
Gosafetmp implements multiple strategies to ensure temporary directories are deleted in most cases, including program halts or (in some cases) even power loss.

## signals catching
Gosafetmp intercepts OS signals to handle process termination. If you wish to handle these signals yourself, set the `SHOULD_CATCH_SIGNALS` variable to `false` before calling `Setup`.

## reaper
Gosafetmp spawns so called "reaper" process that waits for the original process to finish and deletes any remaining tmp directories. If you don't want to use the reaper, set the `SHOULD_SPAWN_REAPER` variable to `false` before calling `Setup`.

## auto deletion
When running on Windows, gosafetmp instructs the OS to delete tmp directories on the next restart (as a fallback if other strategies fail). To disable this, set the `SHOULD_MARK_FOR_AUTO_DELETE` variable to `false` before calling Setup.

# tmpfs
The best way to ensure files don't survive an OS reboot is to place them in an in-memory filesystem (e.g., tmpfs on Linux). While gosafetmp can't enforce this placement, it provides the `IsInTMPFS` function to check if your files reside in an in-memory FS.
