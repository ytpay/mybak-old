package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/cloudfoundry/bytefmt"

	"github.com/sirupsen/logrus"
)

func fullBackup() {

	startTime := now("2006-01-02 15:04:05")

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
		"--parallel=" + strconv.Itoa(Parallel),
		"--compress",
		"--compress-threads=" + strconv.Itoa(CompressThreads),
		"--user=" + User,
		"--password=" + Password,
		"--host=" + Host,
		"--port=" + strconv.Itoa(Port),
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

	endTime := now("2006-01-02 15:04:05")
	if Report {
		tpl := `%s
Start Time: %s
End time: %s
Backup Size: %s
Backup Path: %s`
		banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
		info, err := os.Stat(backupDir)
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Printf(tpl, banner, startTime, endTime, bytefmt.ByteSize(uint64(info.Size())), backupDir)
	}

}
