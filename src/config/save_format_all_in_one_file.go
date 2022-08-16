package config

import "encoding/json"

const (
	allInOneFilenameField = `filename`
)

type SaveFormatAllInOneFile struct {
	filename string
}

func (f *SaveFormatAllInOneFile) GetFormatName() string {
	return formatAllInOne
}

func (f *SaveFormatAllInOneFile) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		allInOneFilenameField: f.filename,
	}
	return json.Marshal(data)
}

func (f *SaveFormatAllInOneFile) Validate() []error {
	return nil
}
