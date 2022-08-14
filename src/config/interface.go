package config

type ConfigInterface interface {
	GetEnv() string
	GetEmail() string
	GetPort() int
	GetDomains() []string
	GetKeyLength() uint16
	GetCertDaysLeftMin() int
	GetUseStaging() bool
	GetStorage() StorageInterface
}

type StorageInterface interface {
	GetAppPath() string
	GetRoot() string
	GetAccountKey() string
	GetCertificateKey() string
	GetCertificates() CertificatesInterface
	GetAccountKeyFullPath() string
	GetCertificateKeyFullPath() string
	GetCertificatesFullChainFullPath() string
	GetCertificatesSplitCertificateFullPath() string
	GetCertificatesSplitIntermediateFullPattern() string
}

type CertificatesInterface interface {
	GetFullChain() string
	GetSplit() SplitInterface
}

type SplitInterface interface {
	GetCertificate() string
	GetIntermediates() string
}
