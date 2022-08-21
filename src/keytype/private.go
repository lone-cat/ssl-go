package keytype

import (
	"crypto/ecdsa"
	"crypto/rsa"
)

type RSA interface {
	*rsa.PrivateKey
}

type EC interface {
	*ecdsa.PrivateKey
}

type Private interface {
	RSA | EC
}
