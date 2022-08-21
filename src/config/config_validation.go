package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"ssl/common"
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

func (c *Config) validateAccountKeyFilename() (errs []error) {
	if c.AccountKeyFilename == `` {
		errs = append(errs, errors.New(`no account key passed`))
		return
	}

	path := filepath.Dir(c.GetAccountKeyFilename())

	exists, _ := common.DirectoryExists(path)
	if !exists {
		errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
	}

	return
}

func (c *Config) validateSaveFormats() (errs []error) {
	if len(c.SaveFormats) < 1 {
		err := errors.New(`less than 1 format passed`)
		errs = append(errs, err)
		return
	}

	err := c.SaveFormats[0].ValidateMain()
	if err != nil {
		errs = append(errs, err)
	}

	for _, format := range c.SaveFormats {
		if format == nil {
			errs = append(errs, errors.New(`nil format passed`))
			continue
		}
		ers := format.Validate()
		errs = append(errs, ers...)
	}

	return
}
