package main

import (
	"os"
)

type ExitCode uint8

const (
	OK ExitCode = iota
	NO_CHANGE
	ERROR
)

func main() {
	os.Exit(int(wrapper()))
}
