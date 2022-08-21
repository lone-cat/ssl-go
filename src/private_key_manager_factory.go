package main

import (
	"os"
	"ssl/keytype"
	"ssl/managers"
	"ssl/storage"
	"ssl/storage/file"
)

func NewPrivateKeyManager[T keytype.Private](filename string, permissions os.FileMode) (mgr managers.PrivateKey[T], err error) {
	byteStorage, err := file.NewByteFile(filename, permissions)
	if err != nil {
		return
	}

	multiByteStorage, err := storage.NewByteSingleFileAdapter(byteStorage)
	if err != nil {
		return
	}

	pemStorage, err := storage.NewPemMultibyte(multiByteStorage)
	if err != nil {
		return
	}

	mgr, err = managers.NewPrivateKey[T](pemStorage)

	return
}
