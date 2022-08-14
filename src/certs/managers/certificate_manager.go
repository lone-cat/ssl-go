package managers

import (
	"crypto/x509"
	"ssl/certs/validations"
)

const (
	DefaultCertificateFilePermissions = 0644
)

var (
	certChainSeparator = []byte("\n")
)

type CertificateManager interface {
	Set(certs []*x509.Certificate) error
	Get() ([]*x509.Certificate, error)
}

type abstractCertificateManager struct {
	certificates []*x509.Certificate
}

func (m *abstractCertificateManager) getValidationError(certs []*x509.Certificate) (err error) {
	err = validations.GetBasicCertificateChainError(certs)
	if err != nil {
		return
	}

	return validations.GetCertificatesOrderError(certs)
}

func (m *abstractCertificateManager) setCertificatesWithoutDump(certs []*x509.Certificate) (err error) {
	err = m.getValidationError(certs)
	if err != nil {
		return
	}

	m.certificates = make([]*x509.Certificate, len(certs))
	copy(m.certificates, certs)

	return
}

func (m *abstractCertificateManager) load() (certs []*x509.Certificate, err error) {
	panic(`this abstract method should not be used!`)
	return
}

func (m *abstractCertificateManager) save() (err error) {
	panic(`this abstract method should not be used!`)
	return
}
