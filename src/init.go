package main

import (
	"errors"
	"ssl/certs"
	"ssl/config"
	loglib "ssl/logger"
)

var logger LoggerInterface

type LoggerInterface interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Error(err error)
}

func init() {
	logger = loglib.Make(`main`)
	config.SetLogger(loglib.Make(`config`))
	certs.SetLogger(loglib.Make(`certs`))
}

func getConfig(envVarKeyEnvironment, envVarKeyConfigFolder string) (conf config.ConfigInterface, err error) {
	conf, errs := config.Initialize(envVarKeyEnvironment, envVarKeyConfigFolder)

	if len(errs) > 0 {
		err = errors.New(`config is not valid`)
		for _, subErr := range errs {
			logger.Error(subErr)
		}
	}

	return
}
