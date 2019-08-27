package common

import (
	"fmt"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	ApiLogger   *logrus.Logger
	MatchLogger *logrus.Logger
)

func InitLogger() (err error) {
	ApiLogger = logrus.New()
	setLogFormat(ApiLogger, config.GetString(config.LogApiFormat))
	apiOutput := config.GetString(config.LogApiOutput)
	if err = setLogOutput(ApiLogger, apiOutput); err != nil {
		fmt.Println("[API server] Set log output error : ", err)
		return
	} else {
		fmt.Println("[API server] Set log output : ", apiOutput)
	}

	apiLevel := config.GetString(config.LogApiLevel)
	if err = setLogLevel(ApiLogger, apiLevel); err != nil {
		fmt.Println("[API server] Set log level error : ", err)
		return
	} else {
		fmt.Println("[API server] Set log level : ", apiLevel)
	}

	MatchLogger = logrus.New()
	setLogFormat(MatchLogger, config.GetString(config.LogMatchFormat))
	matchOutput := config.GetString(config.LogMatchOutput)
	if err = setLogOutput(ApiLogger, matchOutput); err != nil {
		fmt.Println("[Match server] Set log output error : ", err)
		return
	} else {
		fmt.Println("[Match server] Set log output : ", matchOutput)
	}

	matchLevel := config.GetString(config.LogMatchLevel)
	if err = setLogLevel(MatchLogger, matchLevel); err != nil {
		fmt.Println("[Match server] Set log level error : ", err)
		return
	} else {
		fmt.Println("[Match server] Set log level : ", matchLevel)
	}
	return
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
func setLogOutput(log *logrus.Logger, outString string) error {
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
