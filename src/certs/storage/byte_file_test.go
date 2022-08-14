package storage

import (
	"errors"
	"os"
	"testing"
)

const (
	testByteFilename        = `./test/testfilethree.txt`
	testByteFilenameTmp     = `./test/testfile.tmp`
	testByteFilePermissions = 0666
)

var (
	loadFileStore *byteFile
	saveFileStore *byteFile
)

func init() {
	var err error
	loadFileStore, err = NewByteFile(testByteFilename, testByteFilePermissions)
	if err != nil {
		panic(err)
	}
	err = os.Remove(testByteFilenameTmp)
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		panic(err)
	}
	saveFileStore, err = NewByteFile(testByteFilenameTmp, testByteFilePermissions)
	if err != nil {
		panic(err)
	}
}

func TestByteFile_New(t *testing.T) {
	_, err := NewByteFile(testByteFilename, testByteFilePermissions)
	if err != nil {
		t.Fatal(err)
	}
}

func TestByteFile_Load(t *testing.T) {
	data, err := loadFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	reference, err := os.ReadFile(loadFileStore.filename)
	if err != nil {
		t.Fatal(err)
	}

	if !bytesEqual(data, reference) {
		t.Fatal(errors.New(`data does not match`))
	}
}

func TestByteFile_Save(t *testing.T) {
	_, err := saveFileStore.Load()
	if err != NoData {
		t.Fatal(`test file exists`)
	}

	reference := []byte(`tmp data`)

	err = saveFileStore.Save(reference)
	if err != nil {
		t.Fatal(err)
	}

	data, err := saveFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	if !bytesEqual(data, reference) {
		t.Fatal(errors.New(`saved data does not match`))
	}

	reference = []byte(`overwritten tmp data`)
	err = saveFileStore.Save(reference)
	if err != nil {
		t.Fatal(err)
	}

	data, err = saveFileStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	if !bytesEqual(data, reference) {
		t.Fatal(errors.New(`saved data does not match`))
	}
}

func TestByteFile_Delete(t *testing.T) {
	err := saveFileStore.Save([]byte(`some data`))
	if err != nil {
		t.Fatal(err)
	}
	err = saveFileStore.Delete()
	if err != nil {
		t.Fatal(err)
	}
	if fileExists(saveFileStore.filename) {
		t.Fatal(`file was not deleted`)
	}
}
