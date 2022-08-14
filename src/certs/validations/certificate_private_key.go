package validations

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
)

func GetPrivateKeyMatchCertificateError(certificate *x509.Certificate, key *rsa.PrivateKey) error {
	pubKeyFromCert, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New(`public key has improper type`)
	}

	pubKey := &key.PublicKey

	if !pubKeyFromCert.Equal(pubKey) {
		return errors.New(`public key does not match`)
	}

	return nil
}
