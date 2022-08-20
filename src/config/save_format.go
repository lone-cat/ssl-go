package config

const (
	folderField                           = `folder`
	allInOneFilenameField                 = `allInOneFilename`
	privateKeyFilenameField               = `privateKeyFilename`
	certificateFilenameField              = `certificateFilename`
	privateKeyAndCertificateFilenameField = `privateKeyAndCertificateFilename`
	intermediateFilenameField             = `intermediateFilename`
	intermediatePatternField              = `intermediatePattern`
	certificateChainFilenameField         = `certificateChainFilename`
)

type SaveFormat interface {
	Validate() []error
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

func (s *saveFormat) Validate() []error {
	return nil
}
