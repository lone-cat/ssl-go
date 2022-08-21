package main

import (
	"crypto/x509"
	"errors"
	"ssl/config"
	"ssl/keytype"
	"ssl/managers"
)

type MultiBundleManager[T keytype.Private] struct {
	bundleManagers []managers.Bundle[T]
}

func GenerateMultiBundleManagerFromFormatsSlice[T keytype.RSA](saveFormats []config.SaveFormat) (mgr *MultiBundleManager[T], err error) {
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
