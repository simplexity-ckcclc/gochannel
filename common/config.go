package common

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
)

var Conf ConfYaml
var defaultConf = []byte(`
core:
  dsn: ckcclc:141421@tcp(localhost:3306)/gochannel
api:
  address: :8480 
  internal:
    publicKey: publicKey-example
kafka:
  consumer: 
    bootstrap-server: 
      - localhost:9092
    topic: 
      - gochannel
    group-id: gochannel_group

log:
  # text or json
  format: "text"
  # stdout: output to console, or define log path like "log/access_log"
  apiLog: "stdout"
  apiLevel: "debug"

`)

// ConfYaml is config structure.
type ConfYaml struct {
	Core  Core  `yaml:"core"`
	Api   api   `yaml:"api"`
	Kafka kafka `yaml:"kafka"`
	Log   log   `yaml:"log"`
}

type Core struct {
	Database database `yaml:"database"`
}

type database struct {
	DSN          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	//ConnMaxLifetime int `yaml:"conn-max-lifetime"`
}

// SectionCore is sub section of config.
type api struct {
	Address  string   `yaml:"address"`
	Internal internal `yaml:"internal"`
}

type internal struct {
	PublicKey string `yaml:"publicKey"`
}

type kafka struct {
	Consumer consumer `yaml:"consumer"`
}

type consumer struct {
	BootstrapServer []string `yaml:"bootstrap-server"`
	Topic           []string `yaml:"topic"`
	GroupId         string   `yaml:"group-id"`
}

type log struct {
	Format   string `yaml:"format"`
	ApiLog   string `yaml:"apiLog"`
	ApiLevel string `yaml:"apiLevel"`
}

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
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("[Global] Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return err
			}
		}
	}

	// Core
	Conf.Core.Database.DSN = viper.GetString("core.database.dsn")
	Conf.Core.Database.MaxOpenConns = viper.GetInt("core.database.max-open-conns")
	Conf.Core.Database.MaxIdleConns = viper.GetInt("core.database.max-idle-conns")

	// API
	Conf.Api.Address = viper.GetString("api.address")
	Conf.Api.Internal.PublicKey = viper.GetString("api.internal.publicKey")

	// Kafka
	Conf.Kafka.Consumer.BootstrapServer = viper.GetStringSlice("kafka.consumer.bootstrap-server")
	Conf.Kafka.Consumer.Topic = viper.GetStringSlice("kafka.consumer.topic")
	Conf.Kafka.Consumer.GroupId = viper.GetString("kafka.consumer.group-id")

	// Log
	Conf.Log.Format = viper.GetString("log.format")
	Conf.Log.ApiLog = viper.GetString("log.apiLog")
	Conf.Log.ApiLevel = viper.GetString("log.apiLevel")

	return nil
}
