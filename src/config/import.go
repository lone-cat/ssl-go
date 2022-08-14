package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"ssl/common"
	"strings"
)

type EnvStore interface {
	GetEnv() string
	SetEnv(string)
	GetAppPath() string
	SetAppPath(string)
}

func importConfig(config EnvStore, configPath string) error {
	appPath := config.GetAppPath()

	exists, err := common.DirectoryExists(configPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(`config path "` + configPath + `" does not exist`)
	}

	logger.Infof(`importing configuration...`)

	env := config.GetEnv()

	atLeastOneFileImported := false

	configFileNames := generateBasicConfigFileNames()
	for _, configFileName := range configFileNames {
		imported, err := importFileToConfig(config, filepath.Join(configPath, configFileName))
		if err != nil {
			return err
		}
		atLeastOneFileImported = atLeastOneFileImported || imported
	}

	if env != config.GetEnv() {
		if env == `` {
			env = config.GetEnv()
			logger.Infof(`environment set to "%s"`, env)
		} else {
			config.SetEnv(env)
		}
	}

	if env == `` {
		return errors.New(`environment is not set`)
	}

	configFileNames = generateEnvConfigFileNames(env)
	for _, configFileName := range configFileNames {
		imported, err := importFileToConfig(config, filepath.Join(configPath, configFileName))
		if err != nil {
			return err
		}
		atLeastOneFileImported = atLeastOneFileImported || imported
	}

	if !atLeastOneFileImported {
		return errors.New(`no config file found`)
	}

	if env != config.GetEnv() {
		config.SetEnv(env)
	}

	if appPath != config.GetAppPath() {
		config.SetAppPath(appPath)
	}

	logger.Infof(`configuration imported successfully`)

	return nil
}

func generateBasicConfigFileNames() [2]string {
	var filenames [2]string

	point := `.`
	config := `config`
	local := `local`
	extension := `json`

	components := []string{config, extension}
	filenames[0] = strings.Join(components, point)
	components = append(components[:1], local, components[1])
	filenames[1] = strings.Join(components, point)

	return filenames
}

func generateEnvConfigFileNames(env string) [2]string {
	var filenames [2]string

	point := `.`
	config := `config`
	local := `local`
	extension := `json`

	components := []string{config, env, extension}
	filenames[0] = strings.Join(components, point)
	components = append(components[:2], local, components[2])
	filenames[1] = strings.Join(components, point)

	return filenames
}

func importFileToConfig(config EnvStore, path string) (imported bool, err error) {
	imported = false
	var raw []byte

	raw, err = os.ReadFile(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
			logger.Infof(`file "%s" does not exist`, path)
		}
		return
	}

	err = json.Unmarshal(raw, config)
	if err != nil {
		return
	}

	imported = true
	logger.Infof(`file "%s" successfully merged`, path)
	return
}
