package managers

import (
	"encoding/pem"
	"errors"
	"ssl/converters"
	"ssl/keytype"
	"ssl/storage"
)

type PrivateKey[T keytype.Private] interface {
	Get() (T, error)
	Set(T) error
}

type privateKey[T keytype.Private] struct {
	store storage.Pem
}

func NewPrivateKey[T keytype.Private](store storage.Pem) (mgr *privateKey[T], err error) {
	if store == nil {
		err = errors.New(`nil pem storage passed`)
		return
	}

	mgr = &privateKey[T]{
		store: store,
	}

	return
}

func (m *privateKey[T]) Get() (key T, err error) {
	pemBlocks, err := m.store.Load()
	if err != nil {
		return
	}

	keys, _ := converters.PEMBlocksToPrivateKeys[T](pemBlocks)
	if len(keys) < 1 {
		err = errors.New(`no private keys found`)
		return
	}
	key = keys[0]

	return
}

func (m *privateKey[T]) Set(key T) (err error) {
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
