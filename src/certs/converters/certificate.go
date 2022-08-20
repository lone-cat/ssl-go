package converters

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func CertificateToPEMBlock(certificate *x509.Certificate) (pemBlock *pem.Block, err error) {
	if certificate == nil {
		err = errors.New(`nil certificate passed`)
		return
	}

	pemBlock = &pem.Block{Type: `CERTIFICATE`, Bytes: certificate.Raw}

	return
}

func CertificatesToPEMBlocks(certificates []*x509.Certificate) (pemBlocks []*pem.Block, errs []error) {
	pemBlocks = make([]*pem.Block, 0)
	var err error
	if certificates == nil {
		err = errors.New(`nil certificates slice passed`)
		errs = append(errs, err)
		return
	}

	var pemBlock *pem.Block
	for _, certificate := range certificates {
		pemBlock, err = CertificateToPEMBlock(certificate)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		pemBlocks = append(pemBlocks, pemBlock)
	}

	return
}

func PEMBlockToCertificate(pemBlock *pem.Block) (certificate *x509.Certificate, err error) {
	if pemBlock.Type != `CERTIFICATE` {
		err = errors.New(`not certificate block found`)
		return
	}

	certificate, err = x509.ParseCertificate(pemBlock.Bytes)

	return
}

func PEMBlocksToCertificates(pemBlocks []*pem.Block) (certificates []*x509.Certificate, errs []error) {
	certificates = make([]*x509.Certificate, 0)
	var err error
	if pemBlocks == nil {
		err = errors.New(`nil certificates slice passed`)
		errs = append(errs, err)
		return
	}

	var certificate *x509.Certificate
	for _, pemBlock := range pemBlocks {
		certificate, err = PEMBlockToCertificate(pemBlock)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		certificates = append(certificates, certificate)
	}

	return
}
