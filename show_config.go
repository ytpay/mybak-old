package main

import (
	"errors"
	"fmt"
	"strings"

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
		fmt.Printf("mysql user: %s\n", User)
		fmt.Printf("mysql password: %s\n", Password)
		fmt.Printf("mysql host: %s\n", Host)
		fmt.Printf("mysql port: %s\n", Port)
	}
}
