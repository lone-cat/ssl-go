package config

import "errors"

type Certificates struct {
	FullChain string `json:"fullChain"`
	Split     *Split `json:"split"`
}

func (c *Certificates) GetFullChain() string {
	return c.FullChain
}

func (c *Certificates) GetSplit() SplitInterface {
	return c.Split
}

func (c *Certificates) Validate() (errs []error) {
	if c.FullChain == `` && c.Split == nil {
		errs = append(errs, errors.New(`full chain filename should be specified or splitted or both`))
		return
	}
	if c.FullChain != `` {
		errs = append(errs, c.validateFullChain()...)
	}
	if c.Split != nil {
		errs = append(errs, c.validateSplit()...)
	}
	return
}

func (c *Certificates) validateFullChain() (errs []error) {
	return
}

func (c *Certificates) validateSplit() (errs []error) {
	if c.Split != nil {
		errs = append(errs, c.Split.Validate()...)
	}

	return
}
