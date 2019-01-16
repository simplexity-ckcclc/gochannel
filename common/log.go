package common

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
)

func InitLog(conf ConfYaml) error {

	var err error

	// init service logger
	ApiLog = logrus.New()
	setLogFormat(ApiLog, conf.Log.Format)
	if err = setLogOut(ApiLog, conf.Log.ApiLog); err != nil {
		return errors.New("Set api log path error: " + err.Error())
	}

	if err = setLogLevel(ApiLog, conf.Log.ApiLevel); err != nil {
		return errors.New("Set api log level error: " + err.Error())
	}

	return nil
}

func setLogFormat(log *logrus.Logger, format string) {
	switch format {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006/01/02 - 15:04:05",
			FullTimestamp:   true,
		})
	}
}

// SetLogOut provide log stdout and stderr output
func setLogOut(log *logrus.Logger, outString string) error {
	switch outString {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		log.SetOutput(f)
	}
	return nil
}

// SetLogLevel is define log level what you want
// log level: panic, fatal, error, warn, info and debug
func setLogLevel(log *logrus.Logger, levelString string) error {
	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	return nil
}
