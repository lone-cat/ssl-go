package managers

import (
	"crypto"
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
