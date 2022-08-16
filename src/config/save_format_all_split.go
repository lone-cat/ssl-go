package config

import "encoding/json"

const (
	allSplitPrivateKeyField   = `privateKeyFilename`
	allSplitCertificateField  = `certificateFilename`
	allSplitIntermediateField = `intermediatePattern`
)

type SaveFormatAllSplit struct {
	folder              string
	privateKeyFilename  string
	certificateFilename string
	intermediatePattern string
}

func (f *SaveFormatAllSplit) GetFormatName() string {
	return formatAllSplit
}

func (f *SaveFormatAllSplit) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		folderField:               f.folder,
		allSplitPrivateKeyField:   f.privateKeyFilename,
		allSplitCertificateField:  f.certificateFilename,
		allSplitIntermediateField: f.intermediatePattern,
	}
	return json.Marshal(data)
}

func (f *SaveFormatAllSplit) Validate() []error {
	return nil
}
