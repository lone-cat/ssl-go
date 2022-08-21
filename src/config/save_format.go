package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"ssl/common"
)

const (
	defaultPrivateKeyPermissions  = 0600
	defaultCertificatePermissions = 0644
)

type SaveFormat interface {
	Validate() []error
	ValidateMain() error
	GetAllInOneFilename() string
	GetAllInOnePermissions() os.FileMode
	GetPrivateKeyFilename() string
	GetPrivateKeyPermissions() os.FileMode
	GetPrivateKeyAndCertificateFilename() string
	GetPrivateKeyAndCertificatePermissions() os.FileMode
	GetCertificateFilename() string
	GetCertificatePermissions() os.FileMode
	GetCertificateChainFilename() string
	GetCertificateChainPermissions() os.FileMode
	GetIntermediateFilename() string
	GetIntermediatePermissions() os.FileMode
	GetIntermediatePattern() string
	GetIntermediatePatternPermissions() os.FileMode
}

type saveFormat struct {
	Folder                           string `json:"folder"`
	AllInOneFilename                 string `json:"allInOneFilename"`
	PrivateKeyFilename               string `json:"privateKeyFilename"`
	CertificateFilename              string `json:"certificateFilename"`
	PrivateKeyAndCertificateFilename string `json:"privateKeyAndCertificateFilename"`
	IntermediateFilename             string `json:"intermediateFilename"`
	IntermediatePattern              string `json:"intermediatePattern"`
	CertificateChainFilename         string `json:"certificateChainFilename"`
}

func (s *saveFormat) GetAllInOneFilename() string {
	return GenerateFullFilename(s.Folder, s.AllInOneFilename)
}

func (s *saveFormat) GetAllInOnePermissions() os.FileMode {
	return defaultPrivateKeyPermissions
}

func (s *saveFormat) GetPrivateKeyFilename() string {
	return GenerateFullFilename(s.Folder, s.PrivateKeyFilename)
}

func (s *saveFormat) GetPrivateKeyPermissions() os.FileMode {
	return defaultPrivateKeyPermissions
}

func (s *saveFormat) GetPrivateKeyAndCertificateFilename() string {
	return GenerateFullFilename(s.Folder, s.PrivateKeyAndCertificateFilename)
}

func (s *saveFormat) GetPrivateKeyAndCertificatePermissions() os.FileMode {
	return defaultPrivateKeyPermissions
}

func (s *saveFormat) GetCertificateFilename() string {
	return GenerateFullFilename(s.Folder, s.CertificateFilename)
}

func (s *saveFormat) GetCertificatePermissions() os.FileMode {
	return defaultCertificatePermissions
}

func (s *saveFormat) GetCertificateChainFilename() string {
	return GenerateFullFilename(s.Folder, s.CertificateChainFilename)
}

func (s *saveFormat) GetCertificateChainPermissions() os.FileMode {
	return defaultCertificatePermissions
}

func (s *saveFormat) GetIntermediateFilename() string {
	return GenerateFullFilename(s.Folder, s.IntermediateFilename)
}

func (s *saveFormat) GetIntermediatePermissions() os.FileMode {
	return defaultCertificatePermissions
}

func (s *saveFormat) GetIntermediatePattern() string {
	return GenerateFullFilename(s.Folder, s.IntermediatePattern)
}

func (s *saveFormat) GetIntermediatePatternPermissions() os.FileMode {
	return defaultCertificatePermissions
}

func (s *saveFormat) Validate() (errs []error) {
	path := filepath.Dir(s.GetPrivateKeyFilename())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}
	path = filepath.Dir(s.GetCertificateFilename())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}
	path = filepath.Dir(s.GetPrivateKeyAndCertificateFilename())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}
	path = filepath.Dir(s.GetIntermediateFilename())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}
	path = filepath.Dir(s.GetIntermediatePattern())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}
	path = filepath.Dir(s.GetCertificateChainFilename())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}
	path = filepath.Dir(s.GetAllInOneFilename())
	if path != `` {
		exists, _ := common.DirectoryExists(path)
		if !exists {
			errs = append(errs, errors.New(fmt.Sprintf(`folder "%s" does not exist`, path)))
		}
	}

	return
}

func (s *saveFormat) ValidateMain() (err error) {
	if s.GetAllInOneFilename() != `` {
		return
	}

	if s.GetCertificateChainFilename() != `` {
		if s.GetPrivateKeyFilename() != `` {
			return
		}
		if s.GetPrivateKeyAndCertificateFilename() != `` {
			return
		}
	}

	if s.GetIntermediateFilename() != `` || s.GetIntermediatePattern() != `` {
		if s.GetPrivateKeyAndCertificateFilename() != `` {
			return
		}
		if s.GetPrivateKeyFilename() != `` && s.GetCertificateFilename() != `` {
			return
		}
	}

	err = errors.New(`main save format does not contain all necessary data`)
	return
}

func GenerateFullFilename(folder string, filename string) string {
	if filename == `` {
		return ``
	}

	if filepath.IsAbs(filename) {
		return filepath.Clean(filename)
	}

	return filepath.Join(folder, filename)
}
