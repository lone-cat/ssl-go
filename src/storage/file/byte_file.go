package file

import (
	"errors"
	"os"
	"path/filepath"
)

type byteFile struct {
	filename    string
	permissions os.FileMode
}

func NewByteFile(filename string, permissions os.FileMode) (storage *byteFile, err error) {
	folder := filepath.Dir(filename)
	err = validateFolder(folder)
	if err != nil {
		return
	}

	fileBasename := filepath.Base(filename)
	if fileBasename == `` {
		err = errors.New(`empty storage filename passed`)
		return
	}

	storage = &byteFile{
		filename:    filename,
		permissions: permissions,
	}

	return
}

func (s *byteFile) Save(data []byte) error {
	if len(data) < 1 {
		return s.Delete()
	}
	return os.WriteFile(s.filename, data, s.permissions)
}

func (s *byteFile) Load() (bts []byte, err error) {
	bts, err = os.ReadFile(s.filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = NoData
		}
	}

	return
}

func (s *byteFile) Delete() error {
	err := os.Remove(s.filename)
	if os.IsNotExist(err) {
		err = nil
	}
	return err
}
