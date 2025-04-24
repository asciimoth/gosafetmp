package main

import (
	"fmt"
	"os"

	"github.com/asciimoth/gosafetmp"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func pcase(name string) {
	fmt.Println("CASE: " + name)
}

func printMan(man gosafetmp.TmpDirManager) {
	fmt.Println("BASE TMP DIR: ", man.GetBaseDir())
	fmt.Println("IS IN TMPFS: ", man.IsInTMPFS())
}

func printCases(cases map[string](func())) {
	fmt.Println("Available test cases:")
	for k := range cases {
		fmt.Println(k)
	}
}

func tmpDirs(man gosafetmp.TmpDirManager, def bool) {
	for i := 0; i < 10; i++ {
		dir, err := man.NewDir()
		check(err)
		fmt.Println("NEW TMP DIR: ", dir)
		if def {
			defer gosafetmp.Destroy(dir)
		}
	}
}

func dryCase() {
	pcase("Dry run")
	gosafetmp.SHOULD_SPAWN_REAPER = false
	tmpman, err := gosafetmp.Setup()
	check(err)
	printMan(*tmpman)
	tmpDirs(*tmpman, false)
}

func rootDeferCase() {
	pcase("Only root defer")
	gosafetmp.SHOULD_SPAWN_REAPER = false
	tmpman, err := gosafetmp.Setup()
	check(err)
	defer tmpman.Cleanup()
	printMan(*tmpman)
	tmpDirs(*tmpman, false)
}

func tmpsDeferCase() {
	pcase("Only tmp dirs defer")
	gosafetmp.SHOULD_SPAWN_REAPER = false
	tmpman, err := gosafetmp.Setup()
	check(err)
	printMan(*tmpman)
	tmpDirs(*tmpman, true)
}

func reaperCase() {
	pcase("Only reaper guard")
	gosafetmp.SHOULD_SPAWN_REAPER = true
	tmpman, err := gosafetmp.Setup()
	check(err)
	printMan(*tmpman)
	tmpDirs(*tmpman, false)
}

func main() {
	cases := map[string](func()){
		"dry":         dryCase,
		"root-defer":  rootDeferCase,
		"temps-defer": tmpsDeferCase,
		"reaper":      reaperCase,
	}
	if len(os.Args) < 2 {
		printCases(cases)
		return
	}
	cse, ok := cases[os.Args[1]]
	if !ok {
		printCases(cases)
		return
	}
	cse()
}
