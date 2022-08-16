package config

import (
	"encoding/json"
)

type Config struct {
	Env             string   `json:"env"`
	Email           string   `json:"email"`
	Domains         []string `json:"domains"`
	Port            uint16   `json:"port"`
	KeyLength       uint16   `json:"keyLength"`
	CertDaysLeftMin uint8    `json:"certDaysLeftMin"`
	UseStaging      bool     `json:"useStaging"`
	Storage         *Storage `json:"storage"`
	formats         []SaveFormat
	RawFormats      []json.RawMessage `json:"formats"`
}

func NewConfig(env string, appPath string) *Config {
	return &Config{
		Env:        env,
		UseStaging: true,
		Storage:    NewStorage(appPath),
	}
}

func (c *Config) GetEnv() string {
	return c.Env
}

func (c *Config) SetEnv(env string) {
	c.Env = env
}

func (c *Config) GetEmail() string {
	return c.Email
}

func (c *Config) GetPort() int {
	return int(c.Port)
}

func (c *Config) GetDomains() []string {
	domains := make([]string, len(c.Domains))
	copy(domains, c.Domains)
	return domains
}

func (c *Config) GetKeyLength() uint16 {
	return c.KeyLength
}

func (c *Config) GetCertDaysLeftMin() int {
	return int(c.CertDaysLeftMin)
}

func (c *Config) GetUseStaging() bool {
	return c.UseStaging
}

func (c *Config) GetStorage() StorageInterface {
	return c.Storage
}

func (c *Config) GetAppPath() string {
	return c.Storage.AppPath
}

func (c *Config) SetAppPath(appPath string) {
	c.Storage.AppPath = appPath
}

func (c *Config) parseFormats() error {
	c.formats = make([]SaveFormat, 0)
	var err error
	var data interface{}
	var format SaveFormat
	for _, rawJson := range c.RawFormats {
		data = nil
		err = json.Unmarshal(rawJson, &data)
		if err != nil {
			return err
		}
		format, err = convertFormat(data)
		if err != nil {
			return err
		}
		c.formats = append(c.formats, format)
	}
	return nil
}

func (c *Config) String() string {
	startBytes, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}
	var dataMap map[string]interface{}
	err = json.Unmarshal(startBytes, &dataMap)
	if err != nil {
		return err.Error()
	}

	dataMap[`formats`] = c.formats

	jsonBytes, _ := json.MarshalIndent(dataMap, ``, `  `)
	if err != nil {
		return err.Error()
	}

	return string(jsonBytes)
}

func (c *Config) Validate() (errs []error) {
	errs = append(errs, c.validateEnv()...)
	errs = append(errs, c.validateEmail()...)
	errs = append(errs, c.validateDomains()...)
	errs = append(errs, c.validatePort()...)
	errs = append(errs, c.validateKeyLength()...)
	errs = append(errs, c.validateStorage()...)
	return
}
