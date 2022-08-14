package config

import (
	"errors"
	"fmt"
	"regexp"
)

const emailWordSymbols = `[a-z0-9!#$%&'*+/=?^_{|}~-]+`
const emailSymbols = `(?:` + emailWordSymbols + `\.){0,5}` + emailWordSymbols

const domainNameWordSymbols = `[a-z0-9][a-z0-9\-]{0,61}[a-z0-9]`
const domainNameFullRegexp = `(?:` + domainNameWordSymbols + `\.){1,5}` + domainNameWordSymbols

var emailCheckRegexp = regexp.MustCompile(`^` + emailSymbols + `@` + domainNameFullRegexp + `$`)
var domainCheckRegexp = regexp.MustCompile(`^` + domainNameFullRegexp + `$`)

func (c *Config) validateEnv() (errs []error) {
	if c.Env == `` {
		errs = append(errs, errors.New(`environment is not set`))
		return
	}
	if c.Env != `dev` && c.Env != `prod` {
		errs = append(errs, errors.New(fmt.Sprintf(`environment is neither "dev" nor "prod", but "%s"`, c.Env)))
	}
	return
}

func (c *Config) validateEmail() (errs []error) {
	if c.Email == `` {
		errs = append(errs, errors.New(`email is not set`))
		return
	}
	if !emailCheckRegexp.MatchString(c.Email) {
		errs = append(errs, errors.New(`email "`+c.Email+`" is not valid`))
	}
	return
}

func (c *Config) validateDomains() (errs []error) {
	if len(c.Domains) < 1 {
		errs = append(errs, errors.New(`domains are not set`))
		return
	}
	for _, domain := range c.Domains {
		if domain == `` {
			errs = append(errs, errors.New(`domains list contains empty value`))
			continue
		}
		if !domainCheckRegexp.MatchString(domain) {
			errs = append(errs, errors.New(`domain name "`+domain+`" is invalid`))
		}
	}
	return
}

func (c *Config) validatePort() (errs []error) {
	if c.Port < 1 {
		errs = append(errs, errors.New(`port must be in interval 1-65535`))
	}
	return
}

func (c *Config) validateKeyLength() (errs []error) {
	if c.KeyLength < 1 {
		errs = append(errs, errors.New(`key length is not set`))
		return
	}
	if c.KeyLength < 2048 {
		errs = append(errs, errors.New(`key length should not be shorter than 2048 bits`))
	}
	return
}

func (c *Config) validateStorage() (errs []error) {
	if c.Storage == nil {
		errs = append(errs, errors.New(`storage config not filled`))
	} else {
		errs = c.Storage.Validate()
	}
	return
}
