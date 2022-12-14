package config

import (
	"os"
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

	conf.updateFormatFolders()

	logger.Infof("config final version:\n%s", conf)
	errs = conf.Validate()
	if len(errs) < 1 {
		config = conf
		return
	}

	return
}
