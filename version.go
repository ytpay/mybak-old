package main

import (
	"encoding/base64"
	"fmt"
	"runtime"
)

var bannerBase64 = "ICAgIGUgICBlICAgICBZODhiIFk4UCA4ODggODhiLCAgICAgZSBZOGIgICAgIDg4OCA4OFAgCiAgIGQ4YiBkOGIgICAgIFk4OGIgWSAgODg4IDg4UCcgICAgZDhiIFk4YiAgICA4ODggOFAgIAogIGUgWThiIFk4YiAgICAgWTg4YiAgIDg4OCA4SyAgICAgZDg4OGIgWThiICAgODg4IEsgICAKIGQ4YiBZOGIgWThiICAgICA4ODggICA4ODggODhiLCAgZDg4ODg4ODg4OGIgIDg4OCA4YiAgCmQ4ODhiIFk4YiBZOGIgICAgODg4ICAgODg4IDg4UCcgZDg4ODg4ODhiIFk4YiA4ODggODhiIAo="

var versionTpl = `%s
Name: mybak
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
Comment: %s
`

func version() {
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Printf(versionTpl, banner, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID, Comment)
}
