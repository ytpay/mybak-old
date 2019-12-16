package main

import (
	"fmt"
	"runtime"
)

var versionTpl = `
Name: mybak
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
`

func version() {
	fmt.Printf(versionTpl, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID)
}
