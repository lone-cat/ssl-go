package config

import (
	"errors"
	"path/filepath"
	"ssl/common"
)

type Storage struct {
	AppPath        string        `json:"appPath"`
	Root           string        `json:"root"`
	AccountKey     string        `json:"accountKey"`
	CertificateKey string        `json:"certificateKey"`
	Certificates   *Certificates `json:"certificates"`
}

func NewStorage(appPath string) *Storage {
	return &Storage{
		AppPath:      appPath,
		Certificates: &Certificates{},
	}
}

func (s *Storage) GetAppPath() string {
	return s.AppPath
}

func (s *Storage) GetRoot() string {
	return s.Root
}

func (s *Storage) GetAccountKey() string {
	return s.AccountKey
}

func (s *Storage) GetCertificateKey() string {
	return s.CertificateKey
}

func (s *Storage) GetCertificates() CertificatesInterface {
	return s.Certificates
}

func (s *Storage) generateSSLFolderRootPath() string {
	if filepath.IsAbs(s.Root) {
		return s.Root
	}
	return filepath.Join(s.AppPath, s.Root)
}

func (s *Storage) GetAccountKeyFullPath() string {
	if s.AccountKey == `` {
		return ``
	}

	return filepath.Join(s.generateSSLFolderRootPath(), s.AccountKey)
}

func (s *Storage) GetCertificateKeyFullPath() string {
	if s.CertificateKey == `` {
		return ``
	}

	return filepath.Join(s.generateSSLFolderRootPath(), s.CertificateKey)
}

func (s *Storage) GetCertificatesFullChainFullPath() string {
	if s.Certificates == nil {
		return ``
	}
	if s.Certificates.FullChain == `` {
		return ``
	}
	return filepath.Join(s.generateSSLFolderRootPath(), s.Certificates.FullChain)
}

func (s *Storage) GetCertificatesSplitCertificateFullPath() string {
	if s.Certificates == nil {
		return ``
	}
	if s.Certificates.Split == nil {
		return ``
	}
	if s.Certificates.Split.Certificate == `` {
		return ``
	}
	return filepath.Join(s.generateSSLFolderRootPath(), s.Certificates.Split.Certificate)
}

func (s *Storage) GetCertificatesSplitIntermediateFullPattern() string {
	if s.Certificates == nil {
		return ``
	}
	if s.Certificates.Split == nil {
		return ``
	}
	if s.Certificates.Split.Intermediates == `` {
		return ``
	}
	return filepath.Join(s.generateSSLFolderRootPath(), s.Certificates.Split.Intermediates)
}

func (s *Storage) Validate() (errs []error) {
	errs = append(errs, s.validateAppPath()...)
	errs = append(errs, s.validateRoot()...)
	errs = append(errs, s.validateAccountKey()...)
	errs = append(errs, s.validateCertificateKey()...)
	errs = append(errs, s.validateCertificates()...)
	return
}

func (s *Storage) validateAppPath() (errs []error) {
	if !filepath.IsAbs(s.AppPath) {
		errs = append(errs, errors.New(`not absolute application path`))
	}

	exists, err := common.DirectoryExists(s.Root)
	if err != nil {
		errs = append(errs, err)
		return
	}

	if !exists {
		errs = append(errs, errors.New(`app folder does not exist`))
	}

	return
}

func (s *Storage) validateRoot() (errs []error) {
	exists, err := common.DirectoryExists(s.Root)
	if err != nil {
		errs = append(errs, err)
		return
	}

	if !exists {
		errs = append(errs, errors.New(`SSL cert folder does not exist`))
	}

	return
}

func (s *Storage) validateAccountKey() (errs []error) {
	if s.AccountKey == `` {
		errs = append(errs, errors.New(`account key filename not set`))
	}

	return
}

func (s *Storage) validateCertificateKey() (errs []error) {
	if s.CertificateKey == `` {
		errs = append(errs, errors.New(`certificate key filename not set`))
	}

	return
}

func (s *Storage) validateCertificates() (errs []error) {
	if s.Certificates == nil {
		errs = append(errs, errors.New(`no certificates storage data`))
	} else {
		errs = append(errs, s.Certificates.Validate()...)
	}
	return
}
