package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
)

const (
	DatabaseDsn         = "database.dsn"
	DatabaseMaxOpenCons = "database.max-open-conns"
	DatabaseMaxIdleCons = "database.max-idle-conns"

	LogApiFormat   = "logger.api.format"
	LogApiOutput   = "logger.api.output"
	LogApiLevel    = "logger.api.level"
	LogMatchFormat = "logger.match.format"
	LogMatchOutput = "logger.match.output"
	LogMatchLevel  = "logger.match.level"

	EsServer          = "elastic.server"
	EsDeviceIndex     = "elastic.device.index"
	EsDeviceBatchSize = "elastic.device.batch-size"
	EsClickIndex      = "elastic.click.index"
	EsClickBatchSize  = "elastic.click.batch-size"

	ApiServerAddress    = "api.server.address"
	ApiServerIntlPubKey = "api.server.internal.pub-key"

	KafkaBootstrapServer = "kafka.bootstrap-server"
	KafkaTopic           = "kafka.topic"
	KafkaGroupId         = "kafka.group-id"

	ActivateValidPeriod = "activate.valid-period"
	ProcessPeriod       = "process.period"
)

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath string) error {
	viper.SetConfigType("yaml")

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return err
		}
	} else {
		viper.AddConfigPath("$HOME/go")
		viper.AddConfigPath(".")
		viper.SetConfigName("gochannel_config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		fmt.Println("[Global] Using config file:", viper.ConfigFileUsed())
	}
	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}
