package storage

import (
	"encoding/pem"
	"errors"
)

type pemMultibyte struct {
	store ByteMulti
}

func NewPemMultibyte(byteStore ByteMulti) (store *pemMultibyte, err error) {
	if byteStore == nil {
		err = errors.New(`nil byte storage passed`)
		return
	}

	store = &pemMultibyte{
		store: byteStore,
	}

	return
}

func (s *pemMultibyte) Load() (pemBlocks []*pem.Block, err error) {
	data, err := s.store.Load()
	if err != nil {
		return
	}

	if data == nil {
		return
	}

	pemBlocks, err = BytesSliceToPEMBlocks(data)

	return
}

func (s *pemMultibyte) Save(pemBlocks []*pem.Block) (err error) {
	data, err := PEMBlocksToBytesSlice(pemBlocks)
	if err != nil {
		return
	}

	err = s.store.Save(data)

	return
}

func (s *pemMultibyte) Delete() error {
	return s.store.Delete()
}
