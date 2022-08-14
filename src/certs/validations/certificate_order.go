package validations

import "crypto/x509"

func GetCertificatesOrderError(certificateChain []*x509.Certificate) error {
	for i := 0; i < len(certificateChain)-1; i++ {
		cert := certificateChain[i]
		err := cert.CheckSignatureFrom(certificateChain[i+1])
		if err != nil {
			return err
		}
	}

	return nil
}
