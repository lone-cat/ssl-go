package converters

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"go/types"
	"ssl/common"
)

const (
	privateKeyType    = `PRIVATE KEY`
	privateKeyTypeRSA = `RSA PRIVATE KEY`
	privateKeyTypeEC  = `EC PRIVATE KEY`
)

func PrivateKeyToPEMBlock(anyKey crypto.PrivateKey) (pemBlock *pem.Block, err error) {
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

func PrivateKeysToPEMBlocks(anyKeys []crypto.PrivateKey) (pemBlocks []*pem.Block, errs []error) {
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

func PEMBlockToPrivateKey(pemBlock *pem.Block) (key crypto.PrivateKey, err error) {
	switch pemBlock.Type {
	case privateKeyTypeRSA:
		key, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	case privateKeyTypeEC:
		key, err = x509.ParseECPrivateKey(pemBlock.Bytes)
	default:
		key, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	}

	return
}

func PEMBlocksToPrivateKeys(pemBlocks []*pem.Block) (keys []crypto.PrivateKey, errs []error) {
	keys = make([]crypto.PrivateKey, 0)
	var err error
	if pemBlocks == nil {
		err = errors.New(`nil pem blocks slice passed`)
		errs = append(errs, err)
		return
	}

	var key crypto.PrivateKey
	for _, pemBlock := range pemBlocks {
		key, err = PEMBlockToPrivateKey(pemBlock)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		keys = append(keys, key)
	}

	return
}

func ExtractRSAPrivateKeysFromPemBlocks(pemBlocks []*pem.Block) (keys []*rsa.PrivateKey) {
	rawKeys, _ := PEMBlocksToPrivateKeys(pemBlocks)
	keys = common.Convert(rawKeys, rsaConvertFunction)
	return
}

func rsaConvertFunction(rawKey crypto.PrivateKey) (key *rsa.PrivateKey, err error) {
	key, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		err = errors.New(`invalid key type`)
		return
	}

	return
}

func c[T *rsa.PrivateKey|*ecdsa.PrivateKey](a crypto.PrivateKey) T {
	var k T
	k = nil
	fmt.Println(k)
	if types.AssertableTo(a, T)

}