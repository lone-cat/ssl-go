package managers

import (
	"bytes"
	"crypto/x509"
	"errors"
	"ssl/certs/converters"
	"ssl/storage"
	"ssl/storage/file"
)

type certificateChainManager struct {
	*abstractCertificateManager
	certificateStorage storage.Byte
}

func NewCertificateChainManager(certificateStorage storage.Byte) (mgr *certificateChainManager, err error) {
	if certificateStorage == nil {
		err = errors.New(`nil certificate storage passed`)
		return
	}

	mgr = &certificateChainManager{
		abstractCertificateManager: &abstractCertificateManager{},
		certificateStorage:         certificateStorage,
	}

	return
}

func (m *certificateChainManager) Get() (certs []*x509.Certificate, err error) {
	if m.certificates == nil {
		certs, err = m.load()
		if err != nil {
			certs = nil
			return
		}
		m.certificates = certs
	}

	if m.certificates == nil {
		return
	}

	certs = make([]*x509.Certificate, len(m.certificates))
	copy(certs, m.certificates)

	return
}

func (m *certificateChainManager) Set(certs []*x509.Certificate) (err error) {
	err = m.setCertificatesWithoutDump(certs)
	if err != nil {
		return
	}

	return m.save()
}

func (m *certificateChainManager) load() (certs []*x509.Certificate, err error) {
	certBytes, err := m.certificateStorage.Load()
	if err != nil {
		if errors.Is(err, file.NoData) {
			err = nil
		}
		return
	}

	certs, err = converters.ConvertPemBytesToCertificates(certBytes)
	if err != nil {
		return
	}

	return
}

func (m *certificateChainManager) save() (err error) {
	err = m.getValidationError(m.certificates)
	if err != nil {
		return
	}

	certificateBytes, err := converters.ConvertCertificatesToPemBytes(m.certificates)
	if err != nil {
		return err
	}

	return m.certificateStorage.Save(bytes.Join(certificateBytes, certChainSeparator))
}
