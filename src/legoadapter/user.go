package legoadapter

import (
	"crypto"
	"crypto/rsa"
	"github.com/go-acme/lego/v4/registration"
)

// You'll need a user or account type that implements acme.User
type LEUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func GenerateLegoUser(accountPrivateKey *rsa.PrivateKey, email string) *LEUser {
	return &LEUser{
		Email: email,
		key:   accountPrivateKey,
	}
}

func (u *LEUser) GetEmail() string {
	return u.Email
}
func (u *LEUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *LEUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
