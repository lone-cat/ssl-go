package main

import (
	"errors"
)

var NoChangeError = errors.New(`command executed successfully but nothing changed`)

func wrapper() ExitCode {
	logger.Infof(`starting application...`)
	defer logger.Infof(`exited`)

	appConfig, err := getConfig(`APP_ENV`, `APP_CONFIG_FOLDER`)
	if err != nil {
		logger.Error(err)
		return ERROR
	}

	err = app(appConfig)
	if err != nil {
		if err != NoChangeError {
			logger.Error(err)
			return ERROR
		}
		return NO_CHANGE
	}
	return OK
}
