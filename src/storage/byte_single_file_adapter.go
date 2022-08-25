package storage

import (
	"bytes"
	"errors"
)

type byteSingleFileAdapter struct {
	byte      Byte
	separator []byte
}

func NewByteSingleFileAdapter(byteStorage Byte) (store *byteSingleFileAdapter, err error) {
	if byteStorage == nil {
		err = errors.New(`nil byte storage passed`)
		return
	}

	store = &byteSingleFileAdapter{
		byte: byteStorage,
	}

	return
}

func (s *byteSingleFileAdapter) Load() (bts [][]byte, err error) {
	bts = make([][]byte, 0)
	bytesRaw, err := s.byte.Load()
	if err != nil {
		if errors.Is(err, EmptyNode) {
			err = nil
		}
		return
	}

	bts = append(bts, bytesRaw)

	return
}

func (s *byteSingleFileAdapter) Save(bts [][]byte) (err error) {
	return s.byte.Save(bytes.Join(bts, []byte{}))
}

func (s *byteSingleFileAdapter) Delete() error {
	return s.byte.Delete()
}
