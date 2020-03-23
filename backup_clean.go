package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func cleanBackup() {
	matches, err := filepath.Glob(filepath.Join(BackupDir, CleanExpr))
	if err != nil {
		logrus.Fatal(err)
	}

	for _, p := range matches {
		info, err := os.Stat(p)
		if err != nil {
			logrus.Fatal(err)
		}
		if !info.IsDir() {
			continue
		}

		if time.Now().Sub(info.ModTime()) > CleanExpDate {
			if !CleanTest {
				err = os.RemoveAll(p)
				if err != nil {
					logrus.Fatalf("clean backup [%s] failed: %s", p, err)
				}
			}
			logrus.Infof("backup file [%s] has been cleaned.", p)
		}
	}

}
