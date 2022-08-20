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

func Filter[T any](slice []T, filterFunc func(value T) bool) (filteredSlice []T) {
	filteredSlice = make([]T, 0)
	for _, val := range slice {
		if filterFunc(val) {
			filteredSlice = append(filteredSlice, val)
		}
	}

	return
}

func Convert[T any, C any](slice []T, convertFunc func(T) (C, error)) (convertedSlice []C) {
	var err error
	var newVal C
	convertedSlice = make([]C, 0)
	for _, val := range slice {
		newVal, err = convertFunc(val)
		if err != nil {
			continue
		}
		convertedSlice = append(convertedSlice, newVal)
	}

	return
}
