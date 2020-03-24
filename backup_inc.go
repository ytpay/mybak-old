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

func incBackup() {

	startTime := time.Now()

	targetDirTpl, err := template.New("").Funcs(map[string]interface{}{"now": now}).Parse(IncBackupDirTpl)
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
				incBaseDir = strings.TrimSpace(string(bs))
			}
		} else {
			logrus.Fatal(err)
		}
	} else {
		bs, err := ioutil.ReadFile(incBackupStorageFile)
		if err != nil {
			logrus.Fatalf("failed to read inc backup storage file: %v", err)
		}
		incBaseDir = strings.TrimSpace(string(bs))
	}

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
		"--incremental-basedir=" + incBaseDir,
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

	err = ioutil.WriteFile(filepath.Join(BackupDir, IncBackupStorageFile), []byte(backupDir+"\n"), 0644)
	if err != nil {
		logrus.Errorf("failed to storage inc backup dir: %s", backupDir)
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
			err = ioutil.WriteFile(filepath.Join(BackupDir, ReportFile), []byte("```\n"+s+"```"), 0644)
			if err != nil {
				logrus.Fatal(err)
			}
		}
	}
}
