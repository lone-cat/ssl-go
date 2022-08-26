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
	GetCertificate() (*x509.Certificate, error)
	GetIntermediates() ([]*x509.Certificate, error)
	Get() (T, []*x509.Certificate, error)
	Set(T, []*x509.Certificate) error
	NeedSync() bool
	ShouldHavePrivateKey() bool
	ShouldHaveCertificate() bool
	ShouldHaveIntermediates() bool
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
	return m.certificateStorage != nil || m.privateKeyAndCertificateStorage != nil || m.certificateChainStorage != nil || m.allInOneStorage != nil
}

func (m *bundle[T]) ShouldHaveIntermediates() bool {
	return m.certificateChainStorage != nil || m.allInOneStorage != nil || m.intermediateStorage != nil || m.intermediateMultiStorage != nil
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
	var keyBundlesCount int
	if m.privateKeyStorage != nil {
		keyBundlesCount++
	}
	if m.privateKeyAndCertificateStorage != nil {
		keyBundlesCount++
	}
	if m.allInOneStorage != nil {
		keyBundlesCount++
	}

	if keyBundlesCount != len(keyBundles) {
		return true
	}

	certsBundles := m.getCertificatesBundles()
	if certsBundles != nil && len(certsBundles) > 1 {
		certs := certsBundles[0]
		for _, certs2 := range certsBundles[1:] {
			if !CertsBundlesEqual(certs, certs2) {
				return true
			}
		}
	}

	var certsBundlesCount int
	if m.certificateStorage != nil {
		certsBundlesCount++
	}
	if m.privateKeyAndCertificateStorage != nil {
		certsBundlesCount++
	}
	if m.certificateChainStorage != nil {
		certsBundlesCount++
	}
	if m.allInOneStorage != nil {
		certsBundlesCount++
	}

	if certsBundlesCount != len(certsBundles) {
		return true
	}

	certsBundles = m.getIntermediatesBundles()
	if certsBundles != nil && len(certsBundles) > 1 {
		certs := certsBundles[0]
		for _, certs2 := range certsBundles[1:] {
			if !CertsBundlesEqual(certs, certs2) {
				return true
			}
		}
	}

	var intermediateBundlesCount int
	if m.intermediateStorage != nil {
		intermediateBundlesCount++
	}
	if m.intermediateMultiStorage != nil {
		intermediateBundlesCount++
	}
	if m.certificateChainStorage != nil {
		intermediateBundlesCount++
	}
	if m.allInOneStorage != nil {
		intermediateBundlesCount++
	}

	if intermediateBundlesCount != len(certsBundles) {
		return true
	}

	return false
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

func (m *bundle[T]) GetCertificate() (certificate *x509.Certificate, err error) {
	certs := getCertificatesFromPemStorageIfNotNil(m.certificateStorage)
	if len(certs) > 0 {
		certificate = certs[0]
		return
	}

	certs = getCertificatesFromPemStorageIfNotNil(m.privateKeyAndCertificateStorage)
	if len(certs) > 0 {
		certificate = certs[0]
		return
	}

	certs = getCertificatesFromPemStorageIfNotNil(m.certificateChainStorage)
	if len(certs) > 0 {
		certificate = certs[0]
		return

	}

	certs = getCertificatesFromPemStorageIfNotNil(m.allInOneStorage)
	if len(certs) > 0 {
		certificate = certs[0]
	}

	return
}

func (m *bundle[T]) GetIntermediates() (intermediate []*x509.Certificate, err error) {
	certs := getCertificatesFromPemStorageIfNotNil(m.intermediateStorage)
	if certs != nil && (intermediate == nil || len(certs) > len(intermediate)) {
		intermediate = certs
	}

	certs = getCertificatesFromPemStorageIfNotNil(m.intermediateMultiStorage)
	if certs != nil && (intermediate == nil || len(certs) > len(intermediate)) {
		intermediate = certs
	}

	certs = getCertificatesFromPemStorageIfNotNil(m.certificateChainStorage)
	if certs != nil {
		if len(certs) > 0 {
			certs = certs[1:]
		}
		if intermediate == nil || len(certs) > len(intermediate) {
			intermediate = make([]*x509.Certificate, len(certs))
			copy(intermediate, certs)
		}
	}

	certs = getCertificatesFromPemStorageIfNotNil(m.allInOneStorage)
	if certs != nil {
		if len(certs) > 0 {
			certs = certs[1:]
		}
		if intermediate == nil || len(certs) > len(intermediate) {
			intermediate = make([]*x509.Certificate, len(certs))
			copy(intermediate, certs)
		}
	}

	return
}

func (m *bundle[T]) Get() (key T, certificates []*x509.Certificate, err error) {
	key, err = m.GetPrivateKey()
	if err != nil {
		return
	}

	certificate, err := m.GetCertificate()
	if err != nil {
		return
	}

	certificates = append(certificates, certificate)

	intermediate, err := m.GetIntermediates()
	if err != nil {
		return
	}

	certificates = append(certificates, intermediate...)

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

func (m *bundle[T]) getCertificatesBundles() (certificatesBundles [][]*x509.Certificate) {
	var certs, tmp []*x509.Certificate
	certs = getCertificatesFromPemStorageIfNotNil(m.certificateStorage)
	if certs != nil {
		certificatesBundles = append(certificatesBundles, certs)
	}
	certs = getCertificatesFromPemStorageIfNotNil(m.privateKeyAndCertificateStorage)
	if certs != nil {
		certificatesBundles = append(certificatesBundles, certs)
	}
	certs = getCertificatesFromPemStorageIfNotNil(m.certificateChainStorage)
	if certs != nil {
		if len(certs) > 0 {
			tmp = make([]*x509.Certificate, 1)
			copy(tmp, certs)
		} else {
			tmp = make([]*x509.Certificate, 0)
		}
		certificatesBundles = append(certificatesBundles, tmp)
	}
	certs = getCertificatesFromPemStorageIfNotNil(m.allInOneStorage)
	if certs != nil {
		if len(certs) > 0 {
			tmp = make([]*x509.Certificate, 1)
			copy(tmp, certs)
		} else {
			tmp = make([]*x509.Certificate, 0)
		}
		certificatesBundles = append(certificatesBundles, tmp)
	}

	return
}

func (m *bundle[T]) getIntermediatesBundles() (certificatesBundles [][]*x509.Certificate) {
	var certs, tmp []*x509.Certificate
	certs = getCertificatesFromPemStorageIfNotNil(m.intermediateStorage)
	if certs != nil {
		certificatesBundles = append(certificatesBundles, certs)
	}
	certs = getCertificatesFromPemStorageIfNotNil(m.intermediateMultiStorage)
	if certs != nil {
		certificatesBundles = append(certificatesBundles, certs)
	}
	certs = getCertificatesFromPemStorageIfNotNil(m.certificateChainStorage)
	if certs != nil {
		if len(certs) > 0 {
			tmp = make([]*x509.Certificate, len(certs)-1)
			copy(tmp, certs[1:])
		} else {
			tmp = make([]*x509.Certificate, 0)
		}
		certificatesBundles = append(certificatesBundles, tmp)
	}
	certs = getCertificatesFromPemStorageIfNotNil(m.allInOneStorage)
	if certs != nil {
		if len(certs) > 0 {
			tmp = make([]*x509.Certificate, len(certs)-1)
			copy(tmp, certs[1:])
		} else {
			tmp = make([]*x509.Certificate, 0)
		}
		certificatesBundles = append(certificatesBundles, tmp)
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

	if pemBlocks == nil {
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

	if pemBlocks == nil {
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
