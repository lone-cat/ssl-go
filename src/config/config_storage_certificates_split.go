package config

import "errors"

type Split struct {
	Certificate   string `json:"certificate"`
	Intermediates string `json:"intermediates"`
}

func (c *Split) GetCertificate() string {
	return c.Certificate
}

func (c *Split) GetIntermediates() string {
	return c.Intermediates
}

func (c *Split) Validate() (errs []error) {
	errs = append(errs, c.validateCertificate()...)
	errs = append(errs, c.validateIntermediates()...)
	return
}

func (c *Split) validateCertificate() (errs []error) {
	if c.Certificate == `` {
		errs = append(errs, errors.New(`empty certificate filename for split storage`))
	}
	return
}

func (c *Split) validateIntermediates() (errs []error) {
	if c.Intermediates == `` {
		errs = append(errs, errors.New(`empty Intermediates pattern for split storage`))
	}
	return
}
