package storage

import (
	"errors"
)

type byteCacheWrapper struct {
	cachedBytes    []byte
	wrappedStorage Byte
}

func NewByteCacheWrapper(wrapStorage Byte) (*byteCacheWrapper, error) {
	if wrapStorage == nil {
		return nil, errors.New(`nil storage passed`)
	}
	return &byteCacheWrapper{
		wrappedStorage: wrapStorage,
	}, nil
}

func (s *byteCacheWrapper) Save(data []byte) error {
	cache := make([]byte, len(data))
	copy(cache, data)
	s.cachedBytes = cache
	return s.wrappedStorage.Save(cache)
}

func (s *byteCacheWrapper) Load() (data []byte, err error) {
	if s.cachedBytes == nil {
		var bts []byte
		bts, err = s.wrappedStorage.Load()
		if err != nil {
			return
		}
		s.cachedBytes = make([]byte, len(bts))
		copy(s.cachedBytes, bts)
	}

	data = make([]byte, len(s.cachedBytes))
	copy(data, s.cachedBytes)

	return
}

func (s *byteCacheWrapper) Delete() error {
	s.cachedBytes = nil
	return s.wrappedStorage.Delete()
}

func (s *byteCacheWrapper) ClearCache() {
	s.cachedBytes = nil
}
