package validations

import (
	"crypto/x509"
	"errors"
	"fmt"
	"time"
)

const timeFormat = `2006-01-02 15:04:05 MST`

func GetCertificatesExpireError(certificateChain []*x509.Certificate, minLeftTime time.Duration) (err error) {
	curDate := time.Now()
	expireDate := curDate.Add(minLeftTime)

	for _, cert := range certificateChain {
		if curDate.Before(cert.NotBefore) {
			err = errors.New(fmt.Sprintf(`certificate will start at "%s" and now it is "%s"`, cert.NotBefore.In(curDate.Location()).Format(timeFormat), curDate.Format(timeFormat)))
			return
		}
		if expireDate.After(cert.NotAfter) {
			err = errors.New(fmt.Sprintf(`certificate expires at "%s" but it should not expire before "%s"`, cert.NotAfter.In(expireDate.Location()).Format(timeFormat), expireDate.Format(timeFormat)))
			return
		}
	}

	return
}
