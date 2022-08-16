package config

import (
	"errors"
	"os"
	"path/filepath"
	"ssl/common"
)

func Initialize(envVarKeyEnvironment, envVarKeyConfigFolder string) (config ConfigInterface, errs []error) {
	env, exists := os.LookupEnv(envVarKeyEnvironment)
	if exists {
		logger.Infof(`environment from "%s" is set to "%s"`, envVarKeyEnvironment, env)
	} else {
		logger.Infof(`variable "%s" is not set`, envVarKeyEnvironment)
	}

	appPath, err := common.GetAppDir()
	if err != nil {
		errs = append(errs, err)
		return
	}

	conf := NewConfig(env, appPath)

	configPath, exists := os.LookupEnv(envVarKeyConfigFolder)
	if !exists {
		logger.Infof(`config files folder variable "%s" is not set. app root "%s" folder wil be used`, envVarKeyConfigFolder, appPath)
	}

	configPath, err = common.ConvertPathToAbsolute(appPath, configPath)
	if err != nil {
		errs = append(errs, err)
		return
	}

	err = importConfig(conf, configPath)
	if err != nil {
		errs = append(errs, err)
		logger.Infof("config final version:\n%s", conf)
		return
	}

	err = cleanFolders(conf)
	if err != nil {
		errs = append(errs, err)
		logger.Infof("config final version:\n%s", conf)
		return
	}

	err = conf.parseFormats()
	if err != nil {
		errs = append(errs, err)
		logger.Infof("config final version:\n%s", conf)
		return
	}

	logger.Infof("config final version:\n%s", conf)
	errs = conf.Validate()
	if len(errs) < 1 {
		config = conf
		return
	}

	return
}

func cleanFolders(conf *Config) error {
	if conf == nil {
		return errors.New(`nil config passed`)
	}
	if conf.Storage == nil {
		return errors.New(`nil config.Storage`)
	}

	conf.Storage.AppPath = filepath.Clean(conf.Storage.AppPath)
	conf.Storage.Root = filepath.Clean(conf.Storage.Root)

	return nil
}
