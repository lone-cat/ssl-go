package validations

import (
	"crypto/x509"
	"errors"
	"fmt"
)

func GetDomainMatchError(certificate *x509.Certificate, domains []string) error {
	for _, domain := range domains {
		if certificate.VerifyHostname(domain) != nil {
			return errors.New(fmt.Sprintf(`domain "%s" is not included in certificate`, domain))
		}
	}
	return nil
}
