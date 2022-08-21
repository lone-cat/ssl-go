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
	Get() (T, []*x509.Certificate, error)
	Set(T, []*x509.Certificate) error
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

func (m *bundle[T]) Get() (key T, certificates []*x509.Certificate, err error) {
	key, err = m.getPrivateKey()
	if err != nil {
		return
	}

	certificates, err = m.getCertificateChain()
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

	err = m.savePemBlocksToStorageIfNotNil(m.allInOneStorage, pemBlocks)
	if err != nil {
		return
	}

	err = m.savePemBlocksToStorageIfNotNil(m.privateKeyStorage, pemBlocks[:1])
	if err != nil {
		return
	}

	err = m.savePemBlocksToStorageIfNotNil(m.certificateStorage, pemBlocks[1:2])
	if err != nil {
		return
	}

	err = m.savePemBlocksToStorageIfNotNil(m.certificateChainStorage, pemBlocks[1:])
	if err != nil {
		return
	}

	err = m.savePemBlocksToStorageIfNotNil(m.privateKeyAndCertificateStorage, pemBlocks[:2])
	if err != nil {
		return
	}

	err = m.savePemBlocksToStorageIfNotNil(m.intermediateStorage, pemBlocks[2:])
	if err != nil {
		return
	}

	err = m.savePemBlocksToStorageIfNotNil(m.intermediateMultiStorage, pemBlocks[2:])
	if err != nil {
		return
	}

	return
}

func (m *bundle[T]) getPrivateKey() (key T, err error) {
	keys := m.getRSAPrivateKeysFromPemStorage(m.privateKeyStorage)
	if len(keys) > 0 {
		key = keys[0]
		return
	}

	keys = m.getRSAPrivateKeysFromPemStorage(m.privateKeyAndCertificateStorage)
	if len(keys) > 0 {
		key = keys[0]
		return
	}

	keys = m.getRSAPrivateKeysFromPemStorage(m.allInOneStorage)
	if len(keys) > 0 {
		key = keys[0]
		return
	}

	return
}

func (m *bundle[T]) getCertificateChain() (certificates []*x509.Certificate, err error) {
	certChain := m.getCertificatesFromPemStorage(m.certificateChainStorage)
	if len(certChain) > 0 {
		certificates = certChain
		return
	}

	certChain = m.getCertificatesFromPemStorage(m.allInOneStorage)
	if len(certChain) > 0 {
		certificates = certChain
		return
	}

	certs := m.getCertificatesFromPemStorage(m.certificateStorage)
	if len(certs) < 1 {
		certs = m.getCertificatesFromPemStorage(m.privateKeyAndCertificateStorage)
		if len(certs) < 1 {
			return
		}
	}

	intermediate := m.getCertificatesFromPemStorage(m.intermediateStorage)
	if len(intermediate) < 1 {
		intermediate = m.getCertificatesFromPemStorage(m.intermediateMultiStorage)
	}

	certificates = certs
	if len(intermediate) > 0 {
		certificates = append(certificates, intermediate...)
	}

	return
}

func (m *bundle[T]) getRSAPrivateKeysFromPemStorage(store storage.Pem) (keys []T) {
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

func (m *bundle[T]) getCertificatesFromPemStorage(store storage.Pem) (certificates []*x509.Certificate) {
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

func (m *bundle[T]) savePemBlocksToStorageIfNotNil(store storage.Pem, pemBlocks []*pem.Block) (err error) {
	if store != nil {
		err = store.Save(pemBlocks)
	}
	return
}
