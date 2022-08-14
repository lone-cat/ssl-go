package certs

import (
	"crypto/rand"
	"crypto/rsa"
)

func GeneratePrivateKey(len uint16) (*rsa.PrivateKey, error) {
	logger.Infof(`new private key is being generated with length "%d"`, len)
	return rsa.GenerateKey(rand.Reader, int(len))
}
