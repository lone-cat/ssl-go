package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"os"
	"ssl/certs/converters"
	"ssl/storage"
	"ssl/storage/file"
)

type privateKeyManager[T crypto.PrivateKey] struct {
	store storage.Pem
}

func NewPrivateKeyManager[T *ecdsa.PrivateKey | *rsa.PrivateKey](filename string, permissions os.FileMode) (mgr *privateKeyManager[T], err error) {
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

	mgr = &privateKeyManager[T]{
		store: pemStorage,
	}

	return
}

func (m *privateKeyManager[T]) Get() (key T, err error) {
	pemBlocks, err := m.store.Load()
	if err != nil {
		return
	}

	keys, _ := converters.PEMBlocksToPrivateKeys(pemBlocks)
	if len(keys) < 1 {
		err = errors.New(`no private keys found`)
		return
	}
	key = keys[0]

	return
}

func (m *privateKeyManager[T]) GetRSA() (key *rsa.PrivateKey, err error) {
	anyKey, err := m.Get()
	if err != nil {
		return
	}

	key, ok := anyKey.(*rsa.PrivateKey)
	if !ok {
		err = errors.New(`private key is not rsa`)
	}
	return
}

func (m *privateKeyManager) Set(key crypto.PrivateKey) (err error) {
	if key == nil {
		err = errors.New(`nil private key passed`)
		return
	}
	pemBlock, err := converters.PrivateKeyToPEMBlock(key)
	if err != nil {
		return
	}
	err = m.store.Save([]*pem.Block{pemBlock})

	return
}

func (m *privateKeyManager) SetRSA(key *rsa.PrivateKey) error {
	return m.Set(key)
}
