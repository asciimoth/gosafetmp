package main

import (
	"fmt"

	"github.com/asciimoth/gosafetmp"
)

func main() {
	gosafetmp.Setup()
	tmpman, err := gosafetmp.Setup()
	fmt.Println(tmpman, err)
	if err != nil {
		return
	}
	defer tmpman.Cleanup()
	fmt.Println(tmpman.GetBaseDir())
	fmt.Println(tmpman.IsInTMPFS())
	tmpdir, err := tmpman.NewDir()
	fmt.Println(tmpdir, err)
	if err != nil {
		return
	}
	defer gosafetmp.Destroy(tmpdir)
}
