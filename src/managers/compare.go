package managers

import (
	"crypto"
	"crypto/x509"
	"ssl/keytype"
)

func keyBundlesEqual[T keytype.Private](keys1bundle, keys2bundle []T) bool {
	if len(keys1bundle) != len(keys2bundle) {
		return false
	}
	for i := range keys1bundle {
		if !keysEqual(keys1bundle[i], keys2bundle[i]) {
			return false
		}
	}

	return true
}

func keysEqual(key1raw, key2raw crypto.PrivateKey) bool {
	comparableKey1, ok := key1raw.(interface{ Equal(crypto.PrivateKey) bool })
	if !ok {
		return false
	}

	return comparableKey1.Equal(key2raw)
}

func CertsBundlesEqual(certs1bundle, certs2bundle []*x509.Certificate) bool {
	if len(certs1bundle) != len(certs2bundle) {
		return false
	}
	for i := range certs1bundle {
		if !certs1bundle[i].Equal(certs1bundle[i]) {
			return false
		}
	}

	return true
}
