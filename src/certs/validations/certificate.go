package validations

import (
	"crypto/x509"
)

func GetCertificateChainError(certificateChain []*x509.Certificate) (err error) {
	err = GetBasicCertificateChainError(certificateChain)
	if err != nil {
		return
	}

	err = GetCertificatesOrderError(certificateChain)
	if err != nil {
		return
	}

	return
}
