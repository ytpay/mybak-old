package main

import (
	"os"
	"path/filepath"
	"time"
)

var now = func(str string) string {
	return time.Now().Format(str)
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
