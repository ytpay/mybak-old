package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	User                  string
	Password              string
	Threads               int
	Host                  string
	Port                  int
	Secret                string
	Comment               string
	Prefix                string
	BackupDir             string
	FullBackupDirTpl      string
	IncBackupDirTpl       string
	FullBackupStorageFile string
	IncBackupStorageFile  string

	Version   string
	BuildDate string
	CommitID  string
)

var rootCmd = &cobra.Command{
	Use:   "mybak",
	Short: "MySQL backup tool(xtrabackup)",
	Long:  `MySQL backup tool(xtrabackup).`,
	Run:   func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
}

var fullCmd = &cobra.Command{
	Use:   "full",
	Short: "run full backup",
	Long:  `run full backup.`,
	Run:   func(cmd *cobra.Command, args []string) { fullBackup() },
}

var incCmd = &cobra.Command{
	Use:   "inc",
	Short: "run incremental backup",
	Long:  `run incremental backup.`,
	Run:   func(cmd *cobra.Command, args []string) { incBackup() },
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show backup config",
	Long:  `show backup config.`,
	Run:   func(cmd *cobra.Command, args []string) { showConfig() },
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Long:  `show version.`,
	Run:   func(cmd *cobra.Command, args []string) { version() },
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringVar(&Host, "host", "127.0.0.1", "mysql host")
	rootCmd.PersistentFlags().IntVar(&Port, "port", 3306, "mysql port")
	rootCmd.PersistentFlags().IntVar(&Threads, "threads", 8, "backup thread num")
	rootCmd.PersistentFlags().StringVar(&Prefix, "prefix", "mysql", "backup dir prefix")
	rootCmd.PersistentFlags().StringVar(&BackupDir, "backup-dir", "/data/mysql_backup", "backup dir")
	rootCmd.PersistentFlags().StringVar(&FullBackupDirTpl, "full-backup-dir-tpl", `{{ .Prefix }}-{{ "20060102150405" | now }}`, "full backup dir template")
	rootCmd.PersistentFlags().StringVar(&IncBackupDirTpl, "inc-backup-dir-tpl", `{{ .Prefix }}-inc-{{ "20060102150405" | now }}`, "incremental backup dir template")
	rootCmd.PersistentFlags().StringVar(&FullBackupStorageFile, "full-backup-storage-file", ".full-backup", "full backup storage file")
	rootCmd.PersistentFlags().StringVar(&IncBackupStorageFile, "inc-backup-storage-file", ".inc-backup", "incremental backup storage file")
	rootCmd.AddCommand(fullCmd, incCmd, showCmd, versionCmd)
}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
