package managers

import (
	"crypto/x509"
	"errors"
)

type compositeCertificateManager struct {
	mgrs []CertificateManager
}

func NewCompositeCertificateManager(mgrs ...CertificateManager) (man *compositeCertificateManager, err error) {
	if len(mgrs) < 1 {
		err = errors.New(`empty certificate managers slice passed`)
		return
	}

	for _, mgr := range mgrs {
		if mgr == nil {
			err = errors.New(`at least one nil certificate manager passed`)
			return
		}
	}

	realMgrs := make([]CertificateManager, len(mgrs))
	copy(realMgrs, mgrs)

	man = &compositeCertificateManager{
		mgrs: realMgrs,
	}

	return
}

func (m *compositeCertificateManager) Get() ([]*x509.Certificate, error) {
	return m.mgrs[0].Get()
}

func (m *compositeCertificateManager) Set(certs []*x509.Certificate) (err error) {
	for _, mgr := range m.mgrs {
		err = mgr.Set(certs)
		if err != nil {
			return
		}
	}

	return
}
