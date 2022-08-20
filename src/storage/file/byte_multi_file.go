package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const NumberInsertPattern = `{n}`

var patternValidationRegexp = regexp.MustCompile(regexp.QuoteMeta(NumberInsertPattern))

type byteMultiFile struct {
	folder           string
	fileMatchPattern *regexp.Regexp
	fileNameFormat   string
	permissions      os.FileMode
	storages         []*byteFile
}

func NewByteMultiFile(filenamePattern string, permissions os.FileMode) (store *byteMultiFile, err error) {
	folder, fileNameFormat, fileMatchPattern, err := extractPatterns(filenamePattern)
	if err != nil {
		return
	}

	store = &byteMultiFile{
		folder:           folder,
		fileMatchPattern: fileMatchPattern,
		fileNameFormat:   fileNameFormat,
		permissions:      permissions,
		storages:         make([]*byteFile, 0),
	}

	return
}

func (s *byteMultiFile) Load() (bts [][]byte, err error) {
	s.storages, err = getStorageArrayByPattern(s.folder, s.fileMatchPattern, s.permissions)
	if err != nil {
		return
	}

	var fileBytes []byte
	for _, storage := range s.storages {
		fileBytes, err = storage.Load()
		if err != nil {
			return
		}

		if len(fileBytes) > 0 {
			bts = append(bts, fileBytes)
		}
	}

	if bts == nil {
		err = NoData
	}

	return
}

func (s *byteMultiFile) Save(data [][]byte) (err error) {
	err = s.Delete()
	if err != nil {
		return
	}
	maxlen := uint(len(data))
	var filename string
	s.storages = make([]*byteFile, len(data))
	for num, dat := range data {
		filename = s.generateFileNameByIndex(maxlen, uint(num+1))
		s.storages[num], err = NewByteFile(filename, s.permissions)
		if err != nil {
			return
		}
		err = s.storages[num].Save(dat)
		if err != nil {
			return
		}
	}

	return
}

func (s *byteMultiFile) Delete() (err error) {
	s.storages, err = getStorageArrayByPattern(s.folder, s.fileMatchPattern, s.permissions)
	if err != nil {
		return
	}

	for _, storage := range s.storages {
		err = storage.Delete()
		if err != nil {
			return
		}
	}

	s.storages = nil

	return
}

func (s *byteMultiFile) generateFileNameByIndex(maxI, i uint) string {
	return generateFileNameByIndex(s.folder, s.fileNameFormat, maxI, i)
}

func generateFileNameByIndex(folder, format string, maxI, i uint) string {
	maxIstr := strconv.Itoa(int(maxI))
	return filepath.Join(
		folder,
		fmt.Sprintf(format, utf8.RuneCountInString(maxIstr), i),
	)
}

func validateFilenamePattern(filenamePattern string) error {
	patternParts := patternValidationRegexp.Split(filenamePattern, -1)
	if len(patternParts) != 2 {
		return errors.New(`invalid pattern`)
	}

	return nil
}

func extractPatterns(filenamePattern string) (folder, filenameFormat string, fileMatchPattern *regexp.Regexp, err error) {
	folder, filename := filepath.Split(filenamePattern)
	folder = filepath.Clean(folder)

	err = validateFolder(folder)
	if err != nil {
		return
	}

	err = validateFilenamePattern(filename)
	if err != nil {
		return
	}

	filenamePatternParts := patternValidationRegexp.Split(filename, -1)

	filenameFormat = strings.Join(filenamePatternParts, `%0*d`)

	for num := range filenamePatternParts {
		filenamePatternParts[num] = regexp.QuoteMeta(filenamePatternParts[num])
	}

	fileMatchPattern, err = regexp.Compile(strings.Join(filenamePatternParts, `([0-9]+)`))
	if err != nil {
		return
	}

	return
}

func getStorageArrayByPattern(folder string, pattern *regexp.Regexp, permissions os.FileMode) (storages []*byteFile, err error) {
	storages = make([]*byteFile, 0)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}

	var filename string
	var storage *byteFile
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename = file.Name()
		if pattern.MatchString(filename) {
			storage, err = NewByteFile(filepath.Join(folder, filename), permissions)
			if err != nil {
				return
			}
			storages = append(storages, storage)
		}
	}

	return
}
