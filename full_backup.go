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

func fullBackup() {

	targetDirTpl := template.New(FullBackupDirTpl).Funcs(map[string]interface{}{"now": now})
	var buf bytes.Buffer
	err := targetDirTpl.Execute(&buf, struct {
		MySQLName string
	}{
		MySQLName: MySQLName,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	backupDir := filepath.Join(BackupDir, buf.String())

	cmd := exec.Command("xtrabackup",
		"--backup",
		"--dump-innodb-buffer-pool",
		"--compress",
		"--compress-threads=4",
		"--user="+User,
		"--password="+Password,
		"--host="+Host,
		"--port="+Port,
		"--target-dir="+backupDir)

	cmd.Stdout = logrus.StandardLogger().Writer()
	cmd.Stderr = logrus.StandardLogger().Writer()
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(BackupDir, FullBackupStorageFile), []byte(backupDir), 0644)
	if err != nil {
		logrus.Errorf("failed to storage backup dir: %s", backupDir)
		logrus.Fatal(err)
	}

}
