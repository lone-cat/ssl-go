package common

import (
	"os"
	"path/filepath"
)

func GetAppDir() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	path, err = filepath.Abs(filepath.Dir(path))
	if err != nil {
		return "", err
	}

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
