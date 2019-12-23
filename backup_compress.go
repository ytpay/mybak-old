package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/mholt/archiver"
)

func compress(clean bool) {
	var filePaths []string
	paths, err := filepath.Glob(BackupDir + "/*")
	if err != nil {
		logrus.Fatal(err)
	}
	for _, p := range paths {
		_, err := archiver.ByExtension(p)
		if err == nil {
			// skip archive file
			continue
		}
		filePaths = append(filePaths, p)
	}

	logrus.Infof("compress files: %v", filePaths)

	if len(filePaths) < 1 {
		logrus.Fatal("no files to compress")
	}

	bs, err := ioutil.ReadFile(filepath.Join(BackupDir, FullBackupStorageFile))
	if err != nil {
		logrus.Fatal(err)
	}

	targetFile := string(bs) + ".tlz4"
	logrus.Infof("compress target file: %s", targetFile)

	err = archiver.Archive(filePaths, targetFile)
	if err != nil {
		logrus.Fatal(err)
	}
	if clean {
		for _, p := range filePaths {
			err = os.RemoveAll(p)
			if err != nil {
				logrus.Errorf("failed to clean file: %s: %v", p, err)
			}
		}
	}
}

func decompress(src, dist string) {
	err := archiver.Unarchive(src, dist)
	if err != nil {
		logrus.Fatal(err)
	} else {
		logrus.Info("decompress success")
	}
}
