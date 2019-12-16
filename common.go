package main

import "time"

var now = func(str string) string {
	return time.Now().Format(str)
}
