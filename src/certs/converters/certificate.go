package converters

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ConvertCertificatesToPemBytes(certificates []*x509.Certificate) ([][]byte, error) {
	bts := make([][]byte, 0)
	for _, cert := range certificates {
		pemBlock := &pem.Block{Type: `CERTIFICATE`, Bytes: cert.Raw}
		btsBlock := pem.EncodeToMemory(pemBlock)
		if btsBlock == nil {
			return nil, errors.New(`certificate to bytes encoding failed`)
		}
		bts = append(bts, btsBlock)
	}
	return bts, nil
}

func ConvertPemBytesToCertificates(bts []byte) (certs []*x509.Certificate, err error) {
	left := bts
	certs = make([]*x509.Certificate, 0)

	var block *pem.Block
	var cert *x509.Certificate
	for {
		block, left = pem.Decode(left)
		if block == nil {
			break
		}
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return
		}
		if cert == nil {
			err = errors.New(`nil certificate got from bytes`)
			return
		}
		certs = append(certs, cert)
	}

	return
}
