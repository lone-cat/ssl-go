package config

type ConfigInterface interface {
	GetEnv() string
	GetEmail() string
	GetPort() int
	GetDomains() []string
	GetKeyLength() uint16
	GetCertDaysLeftMin() int
	GetUseStaging() bool
	GetAccountKeyFilename() string
	GetSaveFormats() []SaveFormat
	updateFormatFolders()
}
