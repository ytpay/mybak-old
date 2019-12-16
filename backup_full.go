package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

func fullBackup() {

	targetDirTpl, err := template.New("").Funcs(map[string]interface{}{"now": now}).Parse(FullBackupDirTpl)
	if err != nil {
		logrus.Fatal(err)
	}
	var buf bytes.Buffer
	err = targetDirTpl.Execute(&buf, struct {
		Prefix string
	}{
		Prefix: Prefix,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	backupDir := filepath.Join(BackupDir, buf.String())

	cmds := []string{
		"xtrabackup",
		"--backup",
		"--dump-innodb-buffer-pool",
		"--compress",
		"--compress-threads=4",
		"--user=" + User,
		"--password=" + Password,
		"--host=" + Host,
		"--port=" + Port,
		"--target-dir=" + backupDir,
	}

	logrus.Info(strings.Replace(fmt.Sprintf("backup commands: [%s]", strings.Join(cmds, " ")), Password, "********", -1))

	cmd := exec.Command(cmds[0], cmds[1:]...)

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
