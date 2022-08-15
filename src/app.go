package main

import (
	"crypto/rsa"
	"crypto/x509"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"path/filepath"
	"ssl/certs"
	"ssl/certs/converters"
	"ssl/certs/legoadapter"
	"ssl/certs/managers"
	"ssl/certs/storage"
	"ssl/certs/validations"
	"ssl/config"
	"strconv"
	"time"
)

func app(config config.ConfigInterface) error {
	certKeyStorage, err := storage.NewByteFile(config.GetStorage().GetCertificateKeyFullPath(), managers.DefaultPrivateKeyFilePermissions)
	if err != nil {
		return err
	}

	keyManager, err := managers.NewRSAPrivateKeyManager(certKeyStorage)
	if err != nil {
		return err
	}

	var certManagers []managers.CertificateManager
	certificateFullChainFilename := config.GetStorage().GetCertificatesFullChainFullPath()
	if certificateFullChainFilename != `` {
		certificateStorage, err := storage.NewByteFile(certificateFullChainFilename, managers.DefaultCertificateFilePermissions)
		if err != nil {
			return err
		}
		certManager, err := managers.NewCertificateChainManager(certificateStorage)
		if err != nil {
			return err
		}
		certManagers = append(certManagers, certManager)
	}

	splitCertificateFilename := config.GetStorage().GetCertificatesSplitCertificateFullPath()
	intermediateStoragePattern := config.GetStorage().GetCertificatesSplitIntermediateFullPattern()
	if splitCertificateFilename != `` && intermediateStoragePattern != `` {
		certificateStorage, err := storage.NewByteFile(splitCertificateFilename, managers.DefaultCertificateFilePermissions)
		if err != nil {
			return err
		}
		intermediateStorage, err := storage.NewByteMultiFile(intermediateStoragePattern, managers.DefaultCertificateFilePermissions)
		if err != nil {
			return err
		}
		certManager, err := managers.NewCertificateSplitChainManager(certificateStorage, intermediateStorage)
		if err != nil {
			return err
		}
		certManagers = append(certManagers, certManager)
	}

	certManager, err := managers.NewCompositeCertificateManager(certManagers...)
	if err != nil {
		return err
	}

	bundleManager, err := certs.NewBundleManager(keyManager, certManager)
	if err != nil {
		return err
	}

	certKey, certificateChain, err := bundleManager.GetBundle()
	if err != nil {
		return err
	}

	certificateExpireDuration := time.Duration(config.GetCertDaysLeftMin()) * 24 * time.Hour

	if err = validations.GetCertificateBundleValidationError(certKey, certificateChain, config.GetDomains(), certificateExpireDuration); err != nil {
		logger.Error(err)

		certKey, certificateChain, err = getNewCertificateBundle(
			filepath.Join(config.GetStorage().GetAppPath(), config.GetStorage().GetRoot(), config.GetStorage().GetAccountKey()),
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
		err = bundleManager.SetBundle(certKey, certificateChain)
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
	accountKeyStorage, err := storage.NewByteFile(accountKeyFilename, managers.DefaultPrivateKeyFilePermissions)
	if err != nil {
		return
	}

	accountKeyManager, err := managers.NewRSAPrivateKeyManager(accountKeyStorage)
	if err != nil {
		return
	}

	key = accountKeyManager.Get()
	if key == nil || validations.GetRSAPrivateKeyLengthError(key, int(keyLength)) != nil {
		key, err = certs.GeneratePrivateKey(keyLength)
		if err != nil {
			return
		}

		err = accountKeyManager.Set(key)
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

	certificates, err = converters.ConvertPemBytesToCertificates(certificateBytes)
	if err != nil {
		return
	}

	return
}
