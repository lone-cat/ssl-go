package config

import (
	"errors"
	"fmt"
	"path/filepath"
)

const (
	formatAllInOne = `all data in one file`
	formatAllSplit = `all data in separate files`

	folderField = `folder`
)

type SaveFormat interface {
	GetFormatName() string
	Validate() []error
	MarshalJSON() ([]byte, error)
}

func convertFormat(data interface{}) (format SaveFormat, err error) {
	switch data.(type) {
	case string:
		dataMap := map[string]interface{}{allInOneFilenameField: data.(string)}
		format, err = convertMapToFormat(dataMap)
	case map[string]interface{}:
		dataMap := data.(map[string]interface{})
		format, err = convertMapToFormat(dataMap)
	default:
		err = errors.New(`invalid save format entry`)
	}
	return
}

func convertMapToFormat(data map[string]interface{}) (format SaveFormat, err error) {
	if isAllInOneFormat(data) {
		format, err = parseAllInOneFormat(data)
		return
	}
	if isAllSplitFormat(data) {
		format, err = parseAllSplitFormat(data)
		return
	}

	err = errors.New(`format not recognized`)
	return
}

func isAllInOneFormat(data map[string]interface{}) bool {
	_, ok := data[allInOneFilenameField]
	return ok
}

func isAllSplitFormat(data map[string]interface{}) bool {
	match := true
	_, ok := data[allSplitPrivateKeyField]
	match = match && ok
	_, ok = data[allSplitCertificateField]
	match = match && ok
	_, ok = data[allSplitIntermediateField]
	match = match && ok

	return match
}

func parseAllInOneFormat(data map[string]interface{}) (format *SaveFormatAllInOneFile, err error) {
	var folder, filename string
	var ok bool
	folderInterface, folderExist := data[folderField]
	if folderExist {
		folder, ok = tryGetStringFromInterface(folderInterface)
		if !ok {
			err = errors.New(`invalid folder value`)
		}
	}

	filename, ok = tryGetStringFromInterface(data[allInOneFilenameField])
	if !ok {
		fmt.Println(data)
		err = errors.New(`invalid filename value ` + filename)
	}

	if folder != `` {
		filename = filepath.Join(folder, filename)
	}

	format = &SaveFormatAllInOneFile{
		filename: filename,
	}

	return
}

func parseAllSplitFormat(data map[string]interface{}) (format *SaveFormatAllSplit, err error) {
	var folder, privateKeyFilename, certificateFilename, intermediatePattern string
	var ok bool
	folderInterface, folderExist := data[folderField]
	if folderExist {
		folder, ok = tryGetStringFromInterface(folderInterface)
		if !ok {
			err = errors.New(`invalid folder value`)
		}
	}

	privateKeyFilename, ok = tryGetStringFromInterface(data[allSplitPrivateKeyField])
	if !ok {
		err = errors.New(`invalid private key filename value`)
	}

	certificateFilename, ok = tryGetStringFromInterface(data[allSplitCertificateField])
	if !ok {
		err = errors.New(`invalid certificate filename value`)
	}

	intermediatePattern, ok = tryGetStringFromInterface(data[allSplitIntermediateField])
	if !ok {
		err = errors.New(`invalid intermediate pattern value`)
	}

	format = &SaveFormatAllSplit{
		folder:              folder,
		privateKeyFilename:  privateKeyFilename,
		certificateFilename: certificateFilename,
		intermediatePattern: intermediatePattern,
	}

	return
}

func tryGetStringFromInterface(data interface{}) (str string, ok bool) {
	str, ok = data.(string)
	if !ok {
		str = ``
	}

	return
}
