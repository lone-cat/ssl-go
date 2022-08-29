package keytype

import (
	"crypto/ecdsa"
	"crypto/rsa"
)

type Private interface {
	*rsa.PrivateKey | *ecdsa.PrivateKey
}
