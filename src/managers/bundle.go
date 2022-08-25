package managers

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"ssl/converters"
	"ssl/keytype"
	"ssl/storage"
)

type Bundle[T keytype.Private] interface {
	GetPrivateKey() (T, error)
	GetCertificates() ([]*x509.Certificate, error)
	Get() (T, []*x509.Certificate, error)
	Set(T, []*x509.Certificate) error
	NeedSync() bool
	ShouldHavePrivateKey() bool
	ShouldHaveCertificate() bool
}

type bundle[T keytype.Private] struct {
	allInOneStorage                 storage.Pem
	privateKeyStorage               storage.Pem
	certificateStorage              storage.Pem
	privateKeyAndCertificateStorage storage.Pem
	intermediateStorage             storage.Pem
	intermediateMultiStorage        storage.Pem
	certificateChainStorage         storage.Pem
}

func NewBundle[T keytype.Private](
	privateKeyStorage,
	certificateStorage,
	privateKeyAndCertificateStorage,
	certificateChainStorage,
	intermediateStorage,
	intermediateMultiStorage,
	allInOneStorage storage.Pem,
) *bundle[T] {
	return &bundle[T]{
		allInOneStorage:                 allInOneStorage,
		privateKeyStorage:               privateKeyStorage,
		certificateStorage:              certificateStorage,
		privateKeyAndCertificateStorage: privateKeyAndCertificateStorage,
		intermediateStorage:             intermediateStorage,
		intermediateMultiStorage:        intermediateMultiStorage,
		certificateChainStorage:         certificateChainStorage,
	}
}

func (m *bundle[T]) ShouldHavePrivateKey() bool {
	return m.privateKeyStorage != nil || m.privateKeyAndCertificateStorage != nil || m.allInOneStorage != nil
}

func (m *bundle[T]) ShouldHaveCertificate() bool {
	return m.privateKeyStorage != nil || m.privateKeyAndCertificateStorage != nil || m.allInOneStorage != nil
}

func (m *bundle[T]) NeedSync() bool {
	keyBundles := m.getPrivateKeysBundles(m.privateKeyStorage, m.privateKeyAndCertificateStorage, m.allInOneStorage)
	if keyBundles != nil && len(keyBundles) > 1 {
		keys := keyBundles[0]
		for _, keys2 := range keyBundles[1:] {
			if !keyBundlesEqual(keys, keys2) {
				return true
			}
		}
	}

	return false
}

func (m *bundle[T]) Get() (key T, certificates []*x509.Certificate, err error) {
	key, err = m.GetPrivateKey()
	if err != nil {
		return
	}

	certificates, err = m.GetCertificates()
	if err != nil {
		return
	}

	return
}

func (m *bundle[T]) Set(key T, certificates []*x509.Certificate) (err error) {
	keyPemBlock, err := converters.PrivateKeyToPEMBlock(key)
	if err != nil {
		return
	}

	certificateChainBlocks, errs := converters.CertificatesToPEMBlocks(certificates)
	if len(errs) > 0 {
		err = errors.New(`error certificate conversion`)
	}

	pemBlocks := make([]*pem.Block, 0)
	pemBlocks = append(pemBlocks, keyPemBlock)
	pemBlocks = append(pemBlocks, certificateChainBlocks...)

	err = savePemBlocksToStorageIfNotNil(m.allInOneStorage, pemBlocks)
	if err != nil {
		return
	}

	err = savePemBlocksToStorageIfNotNil(m.privateKeyStorage, pemBlocks[:1])
	if err != nil {
		return
	}

	err = savePemBlocksToStorageIfNotNil(m.certificateStorage, pemBlocks[1:2])
	if err != nil {
		return
	}

	err = savePemBlocksToStorageIfNotNil(m.certificateChainStorage, pemBlocks[1:])
	if err != nil {
		return
	}

	err = savePemBlocksToStorageIfNotNil(m.privateKeyAndCertificateStorage, pemBlocks[:2])
	if err != nil {
		return
	}

	err = savePemBlocksToStorageIfNotNil(m.intermediateStorage, pemBlocks[2:])
	if err != nil {
		return
	}

	err = savePemBlocksToStorageIfNotNil(m.intermediateMultiStorage, pemBlocks[2:])
	if err != nil {
		return
	}

	return
}

func (m *bundle[T]) GetPrivateKey() (key T, err error) {
	keys := m.getPrivateKeysFromPemStorage(m.privateKeyStorage)
	if len(keys) > 0 {
		key = keys[0]
		return
	}

	keys = m.getPrivateKeysFromPemStorage(m.privateKeyAndCertificateStorage)
	if len(keys) > 0 {
		key = keys[0]
		return
	}

	keys = m.getPrivateKeysFromPemStorage(m.allInOneStorage)
	if len(keys) > 0 {
		key = keys[0]
		return
	}

	return
}

func (m *bundle[T]) getPrivateKeysBundles(storages ...storage.Pem) (keysBundles [][]T) {
	var keys []T
	for _, store := range storages {
		keys = m.getPrivateKeysFromPemStorage(store)
		if keys != nil {
			keysBundles = append(keysBundles, keys)
		}
	}

	return keysBundles
}

func (m *bundle[T]) GetCertificates() (certificates []*x509.Certificate, err error) {
	certChain := getCertificatesFromPemStorageIfNotNil(m.certificateChainStorage)
	if len(certChain) > 0 {
		certificates = certChain
		return
	}

	certChain = getCertificatesFromPemStorageIfNotNil(m.allInOneStorage)
	if len(certChain) > 0 {
		certificates = certChain
		return
	}

	certs := getCertificatesFromPemStorageIfNotNil(m.certificateStorage)
	if len(certs) < 1 {
		certs = getCertificatesFromPemStorageIfNotNil(m.privateKeyAndCertificateStorage)
		if len(certs) < 1 {
			return
		}
	}

	intermediate := getCertificatesFromPemStorageIfNotNil(m.intermediateStorage)
	if len(intermediate) < 1 {
		intermediate = getCertificatesFromPemStorageIfNotNil(m.intermediateMultiStorage)
	}

	certificates = certs
	if len(intermediate) > 0 {
		certificates = append(certificates, intermediate...)
	}

	return
}

func (m *bundle[T]) getPrivateKeysFromPemStorage(store storage.Pem) (keys []T) {
	if store == nil {
		return
	}

	pemBlocks, err := store.Load()
	if err != nil {
		return
	}

	keys, _ = converters.PEMBlocksToPrivateKeys[T](pemBlocks)

	return
}

func getCertificatesFromPemStorageIfNotNil(store storage.Pem) (certificates []*x509.Certificate) {
	if store == nil {
		return
	}

	pemBlocks, err := store.Load()
	if err != nil {
		return
	}

	certificates, _ = converters.PEMBlocksToCertificates(pemBlocks)

	return
}

func savePemBlocksToStorageIfNotNil(store storage.Pem, pemBlocks []*pem.Block) (err error) {
	if store != nil {
		err = store.Save(pemBlocks)
	}
	return
}
