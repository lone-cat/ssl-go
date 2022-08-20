package file

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	testMultibyteFilename    = `./test/testfile{n}.txt`
	testMultibyteFilenameTmp = `./test/testfile{n}.tmp`
	testMultibytePermissions = 0666
)

var (
	loadMultiFileStore *byteMultiFile
	saveMultiFileStore *byteMultiFile

	tmpMultiReference = [][]byte{
		[]byte(`some data`),
		[]byte(`some more data`),
		[]byte(`some data again`),
	}
)

func init() {
	var err error
	loadMultiFileStore, err = NewByteMultiFile(testMultibyteFilename, testMultibytePermissions)
	if err != nil {
		panic(err)
	}

	saveMultiFileStore, err = NewByteMultiFile(testMultibyteFilenameTmp, testMultibytePermissions)
	if err != nil {
		panic(err)
	}

	err = saveMultiFileStore.Delete()
	if err != nil {
		panic(err)
	}
}

func TestByteMultiFile_generateFilename(t *testing.T) {
	n := 35
	filename := loadMultiFileStore.generateFileNameByIndex(99, uint(n))
	if filename != filepath.Clean(strings.ReplaceAll(testMultibyteFilename, `{n}`, strconv.Itoa(n))) {
		t.Fatal(`filename generate failed`)
	}
}

func TestByteMultiFile_New(t *testing.T) {
	_, err := NewByteMultiFile(testMultibyteFilename, testMultibytePermissions)
	if err != nil {
		t.Fatal(err)
	}
}

func TestByteMultiFile_Load(t *testing.T) {
	data, err := loadMultiFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	data1, err := os.ReadFile(loadMultiFileStore.generateFileNameByIndex(1, 1))
	if err != nil {
		t.Fatal(err)
	}

	data2, err := os.ReadFile(loadMultiFileStore.generateFileNameByIndex(1, 2))
	if err != nil {
		t.Fatal(err)
	}

	reference := [][]byte{data1, data2}
	if !bytesArrEqual(data, reference) {
		t.Fatal(`loaded data does not match`)
	}
}

func TestByteMultiFile_Save(t *testing.T) {
	_, err := saveMultiFileStore.Load()
	if err != NoData {
		t.Fatal(err)
	}

	err = saveMultiFileStore.Save(tmpMultiReference)
	if err != nil {
		t.Fatal(err)
	}

	data, err := saveMultiFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	if !bytesArrEqual(tmpMultiReference, data) {
		t.Fatal(`saved data does not match`)
	}

	newReference := [][]byte{
		[]byte(`changed data`),
		[]byte(`more changed data`),
	}

	err = saveMultiFileStore.Save(newReference)
	if err != nil {
		t.Fatal(err)
	}

	data, err = saveMultiFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	if !bytesArrEqual(newReference, data) {
		t.Fatal(`saved data does not match`)
	}
}

func TestByteMultiFile_Delete(t *testing.T) {
	_, err := saveMultiFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	err = saveMultiFileStore.Delete()
	if err != nil {
		t.Fatal(err)
	}

	_, err = saveMultiFileStore.Load()
	if err != NoData {
		t.Fatal(err)
	}
}
