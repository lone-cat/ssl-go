package validations

import (
	"crypto/x509"
	"errors"
)

func GetBasicCertificateChainError(certificateChain []*x509.Certificate) error {
	if certificateChain == nil {
		return errors.New(`nil certs slice passed`)
	}

	if len(certificateChain) < 1 {
		return errors.New(`empty certs slice passed`)
	}

	for _, cert := range certificateChain {
		if cert == nil {
			return errors.New(`at least one cert is nil in passed slice`)
		}
	}

	return nil
}
