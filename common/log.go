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

func InitLogger() error {
	ApiLogger = logrus.New()
	setLogFormat(ApiLogger, config.GetString(config.LogApiFormat))
	if err := setLogOutput(ApiLogger, config.GetString(config.LogApiOutput)); err != nil {
		fmt.Println("[API server] Set log output error : ", err)
		return err
	}

	if err := setLogLevel(ApiLogger, config.GetString(config.LogApiLevel)); err != nil {
		fmt.Println("[API server] Set log level error : ", err)
		return err
	}

	MatchLogger = logrus.New()
	setLogFormat(MatchLogger, config.GetString(config.LogMatchFormat))
	if err := setLogOutput(ApiLogger, config.GetString(config.LogMatchOutput)); err != nil {
		fmt.Println("[API server] Set log output error : ", err)
		return err
	}

	if err := setLogLevel(MatchLogger, config.GetString(config.LogMatchLevel)); err != nil {
		fmt.Println("[API server] Set log level error : ", err)
		return err
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
	fmt.Println("[API server] Set log output : ", outString)
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
	fmt.Println("[API server] Set log level : ", levelString)
	return nil
}
