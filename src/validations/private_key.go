package validations

import (
	"crypto/rsa"
	"errors"
	"fmt"
)

func GetBasicRSAPrivateKeyError(key *rsa.PrivateKey) error {
	if key == nil {
		return errors.New(`rsa key is nil`)
	}
	return nil
}

func GetRSAPrivateKeyLengthError(key *rsa.PrivateKey, minLen int) error {
	if key.N.BitLen() < minLen {
		return errors.New(fmt.Sprintf(`private key length is less then %d`, minLen))
	}

	return nil
}
