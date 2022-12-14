package main

import (
	"crypto"
	"crypto/x509"
	"errors"
	"ssl/config"
	"ssl/keytype"
	"ssl/managers"
)

type MultiBundleManager[T keytype.Private] struct {
	bundleManagers []managers.Bundle[T]
}

func GenerateMultiBundleManagerFromFormatsSlice[T keytype.Private](saveFormats []config.SaveFormat) (mgr *MultiBundleManager[T], err error) {
	var mgrs []managers.Bundle[T]
	var bundleManager managers.Bundle[T]
	for _, saveFormat := range saveFormats {
		if saveFormat == nil {
			err = errors.New(`at least one format is nil`)
			return
		}

		bundleManager, err = NewBundleManager[T](
			saveFormat.GetPrivateKeyFilename(),
			saveFormat.GetPrivateKeyPermissions(),
			saveFormat.GetCertificateFilename(),
			saveFormat.GetCertificatePermissions(),
			saveFormat.GetPrivateKeyAndCertificateFilename(),
			saveFormat.GetPrivateKeyAndCertificatePermissions(),
			saveFormat.GetCertificateChainFilename(),
			saveFormat.GetCertificateChainPermissions(),
			saveFormat.GetIntermediateFilename(),
			saveFormat.GetIntermediatePermissions(),
			saveFormat.GetIntermediatePattern(),
			saveFormat.GetIntermediatePatternPermissions(),
			saveFormat.GetAllInOneFilename(),
			saveFormat.GetAllInOnePermissions(),
		)
		if err != nil {
			return
		}

		mgrs = append(mgrs, bundleManager)
	}

	return NewMultiBundleManager(mgrs)
}

func NewMultiBundleManager[T keytype.Private](bundleManagers []managers.Bundle[T]) (mgr *MultiBundleManager[T], err error) {
	if len(bundleManagers) < 1 {
		err = errors.New(`empty bundle managers list`)
		return
	}

	for _, bundleManager := range bundleManagers {
		if bundleManager == nil {
			err = errors.New(`at least one bundle manager is nil`)
			return
		}
	}

	mgr = &MultiBundleManager[T]{}

	mgr.bundleManagers = make([]managers.Bundle[T], len(bundleManagers))
	copy(mgr.bundleManagers, bundleManagers)

	return
}

func (m *MultiBundleManager[T]) Sync() (err error) {
	key, certs, _ := m.bundleManagers[0].Get()

	if key == nil || len(certs) < 1 {
		return
	}

	keyComparable, ok := any(key).(interface{ Equal(crypto.PrivateKey) bool })
	if !ok {
		return errors.New(`incomparable key`)
	}

	if m.bundleManagers[0].NeedSync() {
		err = m.bundleManagers[0].Set(key, certs)
		if err != nil {
			return
		}
	}

	for _, mgr := range m.bundleManagers[1:] {
		if mgr.NeedSync() {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

		curKey, _ := mgr.GetPrivateKey()

		if curKey == nil && mgr.ShouldHavePrivateKey() {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

		if curKey != nil && !keyComparable.Equal(curKey) {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

		curCert, _ := mgr.GetCertificate()

		if curCert == nil && mgr.ShouldHaveCertificate() {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

		if curCert != nil && !certs[0].Equal(curCert) {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

		curIntermediates, _ := mgr.GetIntermediates()

		if curIntermediates == nil && mgr.ShouldHaveIntermediates() {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

		if curIntermediates != nil && !managers.CertsBundlesEqual(certs[1:], curIntermediates) {
			err = mgr.Set(key, certs)
			if err != nil {
				return
			}
			continue
		}

	}

	return
}

func (m *MultiBundleManager[T]) GetPrivateKey() (T, error) {
	return m.bundleManagers[0].GetPrivateKey()
}

func (m *MultiBundleManager[T]) GetCertificate() (*x509.Certificate, error) {
	return m.bundleManagers[0].GetCertificate()
}

func (m *MultiBundleManager[T]) GetIntermediates() ([]*x509.Certificate, error) {
	return m.bundleManagers[0].GetIntermediates()
}

func (m *MultiBundleManager[T]) Get() (key T, certificates []*x509.Certificate, err error) {
	return m.bundleManagers[0].Get()
}

func (m *MultiBundleManager[T]) Set(key T, certificates []*x509.Certificate) (err error) {
	for _, mgr := range m.bundleManagers {
		err = mgr.Set(key, certificates)
		if err != nil {
			return
		}
	}

	return
}
