package logger

import (
	"fmt"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"os"
)

type Fields map[string]interface{}

var (
	ApiLogger   *ChannelLogger
	MatchLogger *ChannelLogger
)

// wrap underline logger implementation, only expose generic interface
type ChannelLogger struct {
	*logrus.Logger
	fields map[string]interface{}
}

func newLogger() *ChannelLogger {
	return &ChannelLogger{
		logrus.New(),
		nil,
	}
}

func (logger *ChannelLogger) With(fields map[string]interface{}) *ChannelLogger {
	return &ChannelLogger{
		logger.Logger,
		fields,
	}
}

func (logger *ChannelLogger) Trace(args ...interface{}) {
	logger.Logger.WithFields(logger.fields).Trace(args)
}

func (logger *ChannelLogger) Debug(args ...interface{}) {
	logger.Logger.WithFields(logger.fields).Debug(args)
}

func (logger *ChannelLogger) Info(args ...interface{}) {
	logger.Logger.WithFields(logger.fields).Info(args)
}

func (logger *ChannelLogger) Warn(args ...interface{}) {
	logger.Logger.WithFields(logger.fields).Warn(args)
}

func (logger *ChannelLogger) Error(args ...interface{}) {
	logger.Logger.WithFields(logger.fields).Error(args)
}

func InitLogger() (err error) {
	if err = initApiLogger(); err != nil {
		return
	}

	if err = initMatchLogger(); err != nil {
		return
	}

	return
}

func initApiLogger() (err error) {
	ApiLogger = newLogger()
	ApiLogger.setLogFormat(config.GetString(config.LogApiFormat))
	output := config.GetString(config.LogApiOutput)
	if err = ApiLogger.setLogOutput(output); err != nil {
		fmt.Println("[API server] Set logger output error : ", err)
		return
	} else {
		fmt.Println("[API server] Set logger output : ", output)
	}

	level := config.GetString(config.LogApiLevel)
	if err = ApiLogger.setLogLevel(level); err != nil {
		fmt.Println("[API server] Set logger level error : ", err)
		return
	} else {
		fmt.Println("[API server] Set logger level : ", level)
	}
	return
}

func initMatchLogger() (err error) {
	MatchLogger = newLogger()
	MatchLogger.setLogFormat(config.GetString(config.LogMatchFormat))
	output := config.GetString(config.LogMatchOutput)
	if err = MatchLogger.setLogOutput(output); err != nil {
		fmt.Println("[Match server] Set logger output error : ", err)
		return
	} else {
		fmt.Println("[Match server] Set logger output : ", output)
	}

	level := config.GetString(config.LogMatchLevel)
	if err = MatchLogger.setLogLevel(level); err != nil {
		fmt.Println("[Match server] Set logger level error : ", err)
		return
	} else {
		fmt.Println("[Match server] Set logger level : ", level)
	}
	return
}

func (logger *ChannelLogger) setLogFormat(format string) {
	switch format {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006/01/02 - 15:04:05",
			FullTimestamp:   true,
		})
	}
}

// SetLogOut provide log stdout and stderr output
func (logger *ChannelLogger) setLogOutput(outString string) error {
	switch outString {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		f, err := os.OpenFile(outString, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		logger.SetOutput(f)
	}
	return nil
}

// SetLogLevel is define log level what you want
// log level: panic, fatal, error, warn, info and debug
func (logger *ChannelLogger) setLogLevel(levelString string) error {
	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		return err
	}
	logger.SetLevel(level)
	return nil
}
