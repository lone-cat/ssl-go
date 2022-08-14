package storage

import (
	"errors"
	"os"
)

func fileExists(filename string) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}

	if stat.IsDir() {
		return false
	}

	return true
}

func validateFolder(folder string) (err error) {
	stat, err := os.Stat(folder)
	if err != nil {
		return
	}

	if !stat.IsDir() {
		err = errors.New(`folder does not exists`)
		return
	}

	return
}
