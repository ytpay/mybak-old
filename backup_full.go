package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/cloudfoundry/bytefmt"

	"github.com/sirupsen/logrus"
)

func fullBackup() {

	startTime := time.Now()

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

	err = ioutil.WriteFile(filepath.Join(BackupDir, FullBackupStorageFile), []byte(backupDir+"\n"), 0644)
	if err != nil {
		logrus.Errorf("failed to storage backup dir: %s", backupDir)
		logrus.Fatal(err)
	}

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	if Report {
		tpl := ` __  ____   _____   _   _  __
|  \/  \ \ / / _ ) /_\ | |/ /
| |\/| |\ V /| _ \/ _ \| ' < 
|_|  |_| |_| |___/_/ \_\_|\_\

==============================
Start Time: %s
End Time: %s
Total Time: %s
Backup Size: %s
Backup Path: %s
`
		size, err := DirSize(backupDir)
		if err != nil {
			logrus.Fatal(err)
		}
		s := fmt.Sprintf(tpl, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), totalTime, bytefmt.ByteSize(uint64(size)), backupDir)
		fmt.Println(s)
		if ReportFile != "" {
			err = ioutil.WriteFile(filepath.Join(BackupDir, ReportFile), []byte(s), 0644)
			if err != nil {
				logrus.Fatal(err)
			}
		}
	}

}
