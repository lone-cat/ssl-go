package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"ssl/certs/managers"
	"ssl/certs/validations"
)

type bundleManager struct {
	keyManager         managers.RSAPrivateKeyManager
	certificateManager managers.CertificateManager
}

func NewBundleManager(
	keyManager managers.RSAPrivateKeyManager,
	certificateManager managers.CertificateManager,
) (mgr *bundleManager, err error) {
	if keyManager == nil {
		err = errors.New(`nil key manager passed`)
		return
	}

	if certificateManager == nil {
		err = errors.New(`nil certificate manager passed`)
		return
	}

	mgr = &bundleManager{
		keyManager:         keyManager,
		certificateManager: certificateManager,
	}

	return
}

func (m *bundleManager) GetBundle() (key *rsa.PrivateKey, certs []*x509.Certificate, err error) {
	key = m.keyManager.Get()

	certs, err = m.certificateManager.Get()

	return
}

func (m *bundleManager) SetBundle(key *rsa.PrivateKey, certificates []*x509.Certificate) (err error) {
	err = validations.GetBasicCertificateChainError(certificates)
	if err != nil {
		return
	}

	err = validations.GetCertificatesOrderError(certificates)
	if err != nil {
		return
	}

	err = validations.GetBasicRSAPrivateKeyError(key)
	if err != nil {
		return
	}

	err = validations.GetPrivateKeyMatchCertificateError(certificates[0], key)
	if err != nil {
		return
	}

	err = m.keyManager.Set(key)
	if err != nil {
		return
	}

	return m.certificateManager.Set(certificates)
}
