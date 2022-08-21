package config

import (
	"encoding/json"
	"path/filepath"
)

type Config struct {
	Env             string        `json:"env"`
	Email           string        `json:"email"`
	Domains         []string      `json:"domains"`
	Port            uint16        `json:"port"`
	KeyLength       uint16        `json:"keyLength"`
	CertDaysLeftMin uint8         `json:"certDaysLeftMin"`
	UseStaging      bool          `json:"useStaging"`
	AppPath         string        `json:"appPath"`
	SaveFormats     []*saveFormat `json:"formats"`
}

func NewConfig(env string, appPath string) *Config {
	return &Config{
		Env:        env,
		UseStaging: true,
		AppPath:    appPath,
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

func (c *Config) GetSaveFormats() []SaveFormat {
	if c.SaveFormats == nil {
		return nil
	}

	formats := make([]SaveFormat, len(c.SaveFormats))
	for _, format := range c.SaveFormats {
		formats = append(formats, format)
	}

	return formats
}

func (c *Config) GetAppPath() string {
	return c.AppPath
}

func (c *Config) SetAppPath(appPath string) {
	c.AppPath = appPath
}

func (c *Config) updateFormatFolders() {
	for _, format := range c.SaveFormats {
		if !filepath.IsAbs(format.Folder) {
			format.Folder = filepath.Join(c.AppPath, format.Folder)
		}
	}
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

	dataMap[`SaveFormats`] = c.SaveFormats

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
	errs = append(errs, c.validateSaveFormats()...)
	return
}
