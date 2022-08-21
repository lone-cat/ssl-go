package converters

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"ssl/keytype"
)

const (
	privateKeyType    = `PRIVATE KEY`
	privateKeyTypeRSA = `RSA PRIVATE KEY`
	privateKeyTypeEC  = `EC PRIVATE KEY`
)

func PrivateKeyToPEMBlock[T keytype.Private](anyKey T) (pemBlock *pem.Block, err error) {
	if anyKey == nil {
		err = errors.New(`nil key passed`)
		return
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(anyKey)
	if err != nil {
		return
	}

	pemBlock = &pem.Block{Type: privateKeyType, Bytes: keyBytes}

	return
}

func PrivateKeysToPEMBlocks[T keytype.Private](anyKeys []T) (pemBlocks []*pem.Block, errs []error) {
	pemBlocks = make([]*pem.Block, 0)
	var err error
	if anyKeys == nil {
		err = errors.New(`nil keys slice passed`)
		errs = append(errs, err)
		return
	}

	var pemBlock *pem.Block
	for _, anyKey := range anyKeys {
		pemBlock, err = PrivateKeyToPEMBlock(anyKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		pemBlocks = append(pemBlocks, pemBlock)
	}

	return
}

func PEMBlockToPrivateKey[T keytype.Private](pemBlock *pem.Block) (key T, err error) {
	var anyKey crypto.PrivateKey
	switch pemBlock.Type {
	case privateKeyTypeRSA:
		anyKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	case privateKeyTypeEC:
		anyKey, err = x509.ParseECPrivateKey(pemBlock.Bytes)
	default:
		anyKey, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	}
	if err == nil {
		var ok bool
		key, ok = anyKey.(T)
		if !ok {
			err = errors.New(`invalid key`)
		}
	}

	return
}

func PEMBlocksToPrivateKeys[T keytype.Private](pemBlocks []*pem.Block) (keys []T, errs []error) {
	keys = make([]T, 0)
	var err error
	if pemBlocks == nil {
		err = errors.New(`nil pem blocks slice passed`)
		errs = append(errs, err)
		return
	}

	var key T
	for _, pemBlock := range pemBlocks {
		key, err = PEMBlockToPrivateKey[T](pemBlock)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		keys = append(keys, key)
	}

	return
}
