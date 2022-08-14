package converters

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ConvertPrivateKeyToPemBytes(key *rsa.PrivateKey) (bytes []byte, err error) {
	pkeyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return
	}

	pemBlock := &pem.Block{Type: `PRIVATE KEY`, Bytes: pkeyBytes}
	bytes = pem.EncodeToMemory(pemBlock)

	return
}

func ConvertPemBytesToPrivateKey(bts []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(bts)
	if block == nil {
		err = errors.New(`no pem block found`)
		return
	}
	pKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	switch pKey.(type) {
	case *rsa.PrivateKey:
		key = pKey.(*rsa.PrivateKey)
	default:
		err = errors.New(`invalid key type`)
	}
	return
}
