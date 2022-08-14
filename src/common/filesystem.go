package common

import (
	"errors"
	"os"
	"path/filepath"
)

func EnsureDirectoryExists(path string) error {
	exists, err := DirectoryExists(path)
	if err == nil && !exists {
		err = errors.New(`folder does not exists (` + path + `)`)
	}
	return err
}

func DirectoryExists(path string) (exists bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return
	}
	exists = stat.IsDir()
	return
}

func FileExists(path string) (exists bool, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return
	}
	exists = !stat.IsDir()
	return
}

func ConvertPathToAbsolute(rootPath string, path string) (newPath string, err error) {
	newPath = path
	if filepath.IsAbs(path) {
		return
	}

	newPath, err = filepath.Abs(filepath.Join(rootPath, newPath))
	if err != nil {
		return
	}

	newPath, err = filepath.EvalSymlinks(newPath)
	if err != nil {
		return
	}

	return
}
