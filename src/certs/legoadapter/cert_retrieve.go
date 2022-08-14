package legoadapter

import (
	"crypto/rsa"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
)

func RequestCertificateBytesForDomains(client *lego.Client, domains []string, certPrivateKey *rsa.PrivateKey) (certBytes []byte, err error) {
	request := certificate.ObtainRequest{
		Domains:    domains,
		Bundle:     true,
		PrivateKey: certPrivateKey,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	return certificates.Certificate, nil
}
