package validations

import (
	"crypto/rsa"
	"crypto/x509"
	"time"
)

func GetCertificateBundleValidationError(
	certKey *rsa.PrivateKey,
	certificateChain []*x509.Certificate,
	domains []string,
	minLeftTime time.Duration,
) (err error) {
	err = GetBasicCertificateChainError(certificateChain)
	if err != nil {
		return
	}

	err = GetCertificatesOrderError(certificateChain)
	if err != nil {
		return
	}

	err = GetCertificatesExpireError(certificateChain, minLeftTime)
	if err != nil {
		return
	}

	err = GetBasicRSAPrivateKeyError(certKey)
	if err != nil {
		return
	}

	certificate := certificateChain[0]

	err = GetPrivateKeyMatchCertificateError(certificate, certKey)
	if err != nil {
		return
	}

	err = GetDomainMatchError(certificate, domains)
	if err != nil {
		return
	}

	return
}
