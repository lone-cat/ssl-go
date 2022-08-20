package managers

import (
	"crypto/rsa"
	"errors"
	"ssl/certs/converters"
	"ssl/certs/validations"
	"ssl/storage"
	"ssl/storage/file"
)

const (
	DefaultPrivateKeyFilePermissions = 0600
)

type RSAPrivateKeyManager interface {
	Set(key *rsa.PrivateKey) (err error)
	Get() *rsa.PrivateKey
}

type rsaPrivateKeyManager struct {
	key     *rsa.PrivateKey
	storage storage.Byte
}

func NewRSAPrivateKeyManager(storage storage.Byte) (mgr *rsaPrivateKeyManager, err error) {
	if storage == nil {
		err = errors.New(`nil private key storage passed`)
		return
	}

	mgr = &rsaPrivateKeyManager{
		storage: storage,
	}

	err = mgr.init()

	return
}

func (m *rsaPrivateKeyManager) init() (err error) {
	keyBytes, err := m.storage.Load()
	if err != nil {
		if errors.Is(err, file.NoData) {
			err = nil
		}
		return
	}

	key, err := converters.ConvertPemBytesToPrivateKey(keyBytes)
	if err != nil {
		return
	}

	return m.setKeyWithoutDump(key)
}

func (m *rsaPrivateKeyManager) getValidationError(key *rsa.PrivateKey) (err error) {
	return validations.GetBasicRSAPrivateKeyError(key)
}

func (m *rsaPrivateKeyManager) setKeyWithoutDump(key *rsa.PrivateKey) (err error) {
	err = m.getValidationError(key)
	if err != nil {
		return
	}

	m.key = key

	return
}

func (m *rsaPrivateKeyManager) dumpKey() (err error) {
	err = m.getValidationError(m.key)
	if err != nil {
		return
	}

	keyBytes, err := converters.ConvertPrivateKeyToPemBytes(m.key)
	if err != nil {
		return err
	}

	return m.storage.Save(keyBytes)
}

func (m *rsaPrivateKeyManager) Set(key *rsa.PrivateKey) (err error) {
	err = m.setKeyWithoutDump(key)
	if err != nil {
		return
	}

	return m.dumpKey()
}

func (m *rsaPrivateKeyManager) Get() *rsa.PrivateKey {
	return m.key
}
