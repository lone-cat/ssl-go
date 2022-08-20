package main

import (
	"crypto/rsa"
	"crypto/x509"
	"ssl/certs/converters"
	"ssl/storage"
)

type BundleManager struct {
	allInOneStorage                 storage.Pem
	privateKeyStorage               storage.Pem
	certificateStorage              storage.Pem
	privateKeyAndCertificateStorage storage.Pem
	intermediateStorage             storage.Pem
	intermediateMultiStorage        storage.Pem
	certificateChainStorage         storage.Pem
}

func NewBundleManager(
	allInOneStorage storage.Pem,
	privateKeyStorage storage.Pem,
	certificateStorage storage.Pem,
	privateKeyAndCertificateStorage storage.Pem,
	intermediateStorage storage.Pem,
	intermediateMultiStorage storage.Pem,
	certificateChainStorage storage.Pem,
) *BundleManager {
	return &BundleManager{
		allInOneStorage:                 allInOneStorage,
		privateKeyStorage:               privateKeyStorage,
		certificateStorage:              certificateStorage,
		privateKeyAndCertificateStorage: privateKeyAndCertificateStorage,
		intermediateStorage:             intermediateStorage,
		intermediateMultiStorage:        intermediateMultiStorage,
		certificateChainStorage:         certificateChainStorage,
	}
}

func (m *BundleManager) Get() (key *rsa.PrivateKey, certificates []*x509.Certificate, err error) {
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

func (m *BundleManager) Set(key *rsa.PrivateKey, certificates []*x509.Certificate) error {
	return nil
}

func (m *BundleManager) getPrivateKey() (key *rsa.PrivateKey, err error) {
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

func (m *BundleManager) getCertificateChain() (certificates []*x509.Certificate, err error) {
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

}

func (m *BundleManager) getRSAPrivateKeysFromPemStorage(store storage.Pem) (keys []*rsa.PrivateKey) {
	if store == nil {
		return
	}

	pemBlocks, err := store.Load()
	if err != nil {
		return
	}

	keys = converters.ExtractRSAPrivateKeysFromPemBlocks(pemBlocks)

	return
}

func (m *BundleManager) getCertificatesFromPemStorage(store storage.Pem) (certificates []*x509.Certificate) {
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
