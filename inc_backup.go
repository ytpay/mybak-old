package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/sirupsen/logrus"
)

func incBackup() {
	targetDirTpl := template.New(IncBackupDirTpl).Funcs(map[string]interface{}{"now": now})
	var buf bytes.Buffer
	err := targetDirTpl.Execute(&buf, struct {
		MySQLName string
	}{
		MySQLName: MySQLName,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	var incBaseDir string
	backupDir := filepath.Join(BackupDir, buf.String())
	fullBackupStorageFile := filepath.Join(BackupDir, FullBackupStorageFile)
	incBackupStorageFile := filepath.Join(BackupDir, IncBackupStorageFile)

	_, err = os.Stat(incBackupStorageFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("unable to find incremental backup file: %s", incBackupStorageFile)
			_, err = os.Stat(fullBackupStorageFile)
			if err != nil {
				if os.IsNotExist(err) {
					logrus.Fatalf("unable to find incremental backup file and full backup file, incremental backup failed.")
				} else {
					logrus.Fatal(err)
				}
			} else {
				bs, err := ioutil.ReadFile(fullBackupStorageFile)
				if err != nil {
					logrus.Fatalf("failed to read inc backup storage file: %v", err)
				}
				incBaseDir = string(bs)
			}
		} else {
			logrus.Fatal(err)
		}
	} else {
		bs, err := ioutil.ReadFile(incBackupStorageFile)
		if err != nil {
			logrus.Fatalf("failed to read inc backup storage file: %v", err)
		}
		incBaseDir = string(bs)
	}

	cmd := exec.Command("xtrabackup",
		"--backup",
		"--dump-innodb-buffer-pool",
		"--compress",
		"--compress-threads=4",
		"--user="+User,
		"--password="+Password,
		"--host="+Host,
		"--port="+Port,
		"--incremental-basedir="+incBaseDir,
		"--target-dir="+backupDir)

	cmd.Stdout = logrus.StandardLogger().Writer()
	cmd.Stderr = logrus.StandardLogger().Writer()
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(BackupDir, IncBackupStorageFile), []byte(backupDir), 0644)
	if err != nil {
		logrus.Errorf("failed to storage inc backup dir: %s", backupDir)
		logrus.Fatal(err)
	}
}
