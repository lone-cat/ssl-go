package legoadapter

import (
	"github.com/go-acme/lego/v4/acme"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

func GetLegoClient(user registration.User, useStagingCA bool) (*lego.Client, error) {
	LEconfig := lego.NewConfig(user)

	if useStagingCA {
		LEconfig.CADirURL = lego.LEDirectoryStaging
	} else {
		LEconfig.CADirURL = lego.LEDirectoryProduction
	}

	LEconfig.Certificate.KeyType = certcrypto.RSA4096

	// A client facilitates communication with the CA server.
	return lego.NewClient(LEconfig)
}

func LoginOrRegisterIfNotExists(client *lego.Client) (resource *registration.Resource, err error) {
	resource, err = client.Registration.ResolveAccountByKey()
	if err != nil {
		er, ok := err.(*acme.ProblemDetails)
		if ok && er.Type == `urn:ietf:params:acme:error:accountDoesNotExist` {
			// New users will need to register
			resource, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		}
	}

	return
}
