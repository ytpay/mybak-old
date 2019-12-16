package main

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/mritd/promptx"
)

func showConfig() {
	p := promptx.NewDefaultPrompt(func(line []rune) error {
		if strings.TrimSpace(string(line)) == "" {
			return errors.New("input is empty")
		} else {
			return nil
		}
	}, "Please input secret:")

	if p.Run() == Secret {
		logrus.Infof("mysql user: %s\n", User)
		logrus.Infof("mysql password: %s\n", Password)
		logrus.Infof("mysql host: %s\n", Host)
		logrus.Infof("mysql port: %s\n", Port)
	} else {
		logrus.Fatal("go fuck yourself")
	}
}
