package main

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"ssl/certs"
	"ssl/config"
	"ssl/converters"
	"ssl/legoadapter"
	"ssl/storage"
	"ssl/storage/memory"
	"ssl/validations"
	"strconv"
	"time"
)

func app(config config.ConfigInterface) (err error) {
	bundleManager, err := GenerateMultiBundleManagerFromFormatsSlice[*rsa.PrivateKey](config.GetSaveFormats())
	if err != nil {
		return
	}

	certKey, certificateChain, err := bundleManager.bundleManagers[0].Get()
	if err != nil {
		return
	}

	certificateExpireDuration := time.Duration(config.GetCertDaysLeftMin()) * 24 * time.Hour

	if err = validations.GetCertificateBundleValidationError(certKey, certificateChain, config.GetDomains(), certificateExpireDuration); err != nil {
		logger.Error(err)

		certKey, certificateChain, err = getNewCertificateBundle(
			config.GetAccountKeyFilename(),
			config.GetKeyLength(),
			config.GetEmail(),
			config.GetDomains(),
			config.GetPort(),
			config.GetUseStaging(),
		)
		if err != nil {
			return err
		}

		// TODO: order certificates in chain so cert is first, later trust chain in child-to-parent order
		err = bundleManager.Set(certKey, certificateChain)
		if err != nil {
			return err
		}

		err = validations.GetCertificateBundleValidationError(certKey, certificateChain, config.GetDomains(), certificateExpireDuration)
		if err != nil {
			logger.Errorf(`retrieved certs are invalid: %s`, err.Error())
		}

		return nil
	} else {
		logger.Infof(`certificate bundle is ok`)
		return NoChangeError
	}

}

func getNewCertificateBundle(accountKeyFilename string, keyLength uint16, email string, domains []string, port int, useStagingCA bool) (key *rsa.PrivateKey, certificates []*x509.Certificate, err error) {
	accountKey, err := getOrGenerateAccountKey(accountKeyFilename, keyLength)
	if err != nil {
		return
	}

	client, err := getConnectedClient(accountKey, email, useStagingCA)
	if err != nil {
		return
	}

	key, err = certs.GeneratePrivateKey(keyLength)
	if err != nil {
		return
	}

	certificates, err = getCertificates(client, key, domains, port)

	return
}

func getOrGenerateAccountKey(accountKeyFilename string, keyLength uint16) (key *rsa.PrivateKey, err error) {
	mgr, err := NewPrivateKeyManager[*rsa.PrivateKey](accountKeyFilename, 0600)
	if err != nil {
		return
	}

	key, err = mgr.Get()
	if err != nil {
		if !errors.Is(err, storage.NoData) {
			return
		}
		err = nil
	}

	if key == nil || validations.GetRSAPrivateKeyLengthError(key, int(keyLength)) != nil {
		key, err = certs.GeneratePrivateKey(keyLength)
		if err != nil {
			return
		}

		err = mgr.Set(key)
		if err != nil {
			return
		}
	}

	return
}

func getConnectedClient(accountKey *rsa.PrivateKey, email string, useStagingCA bool) (client *lego.Client, err error) {
	user := legoadapter.GenerateLegoUser(accountKey, email)

	client, err = legoadapter.GetLegoClient(user, useStagingCA)
	if err != nil {
		return
	}

	resource, err := legoadapter.LoginOrRegisterIfNotExists(client)
	if err != nil {
		return
	}
	user.Registration = resource

	return
}

func getCertificates(client *lego.Client, key *rsa.PrivateKey, domains []string, port int) (certificates []*x509.Certificate, err error) {
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer(``, strconv.Itoa(port)))
	if err != nil {
		return
	}

	certificateBytes, err := legoadapter.RequestCertificateBytesForDomains(client, domains, key)
	if err != nil {
		return
	}

	byteStore := memory.NewByteMemory()
	_ = byteStore.Save(certificateBytes)

	multiByteStore, err := storage.NewByteSingleFileAdapter(byteStore)
	if err != nil {
		return
	}

	pemStore, err := storage.NewPemMultibyte(multiByteStore)
	if err != nil {
		return
	}

	pemBlocks, err := pemStore.Load()
	if err != nil {
		return
	}

	certificates, errs := converters.PEMBlocksToCertificates(pemBlocks)
	if len(errs) > 0 {
		err = errs[0]
		return
	}

	return
}
