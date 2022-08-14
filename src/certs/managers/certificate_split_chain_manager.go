package managers

import (
	"bytes"
	"crypto/x509"
	"errors"
	"ssl/certs/converters"
	"ssl/certs/storage"
)

type certificateSplitChainManager struct {
	*abstractCertificateManager
	certificateStorage storage.Byte
	caStorage          storage.ByteMulti
}

func NewCertificateSplitChainManager(
	certificateStorage storage.Byte,
	caStorage storage.ByteMulti,
) (mgr *certificateSplitChainManager, err error) {
	if certificateStorage == nil {
		err = errors.New(`nil certificate storage passed`)
		return
	}
	if caStorage == nil {
		err = errors.New(`nil certificate authority storage passed`)
		return
	}

	mgr = &certificateSplitChainManager{
		abstractCertificateManager: &abstractCertificateManager{},
		certificateStorage:         certificateStorage,
		caStorage:                  caStorage,
	}

	return
}

func (m *certificateSplitChainManager) Get() (certs []*x509.Certificate, err error) {
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

func (m *certificateSplitChainManager) Set(certs []*x509.Certificate) (err error) {
	err = m.setCertificatesWithoutDump(certs)
	if err != nil {
		return
	}

	return m.save()
}

func (m *certificateSplitChainManager) load() (certs []*x509.Certificate, err error) {
	certBytes, err := m.certificateStorage.Load()
	if err != nil {
		if errors.Is(err, storage.NoData) {
			err = nil
		}
		return
	}
	certs, err = converters.ConvertPemBytesToCertificates(certBytes)
	if err != nil {
		return
	}

	intermediateBytes, err := m.caStorage.Load()
	if err != nil {
		if errors.Is(err, storage.NoData) {
			err = nil
		}
		return
	}
	certBytes = bytes.Join(intermediateBytes, certChainSeparator)
	intermediateCerts, err := converters.ConvertPemBytesToCertificates(certBytes)
	if err != nil {
		return
	}

	certs = append(certs, intermediateCerts...)

	return
}

func (m *certificateSplitChainManager) save() (err error) {
	err = m.getValidationError(m.certificates)
	if err != nil {
		return
	}

	certificateBytes, err := converters.ConvertCertificatesToPemBytes(m.certificates[:1])
	if err != nil {
		return
	}

	err = m.certificateStorage.Save(certificateBytes[0])
	if err != nil {
		return
	}

	certificateBytes, err = converters.ConvertCertificatesToPemBytes(m.certificates[1:])
	if err != nil {
		return
	}

	err = m.caStorage.Save(certificateBytes)

	return
}
