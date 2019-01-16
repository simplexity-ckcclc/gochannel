package common

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
)

var defaultConf = []byte(`
core:
  dsn: ckcclc:141421@tcp(localhost:3306)/gochannel
api:
  internal: 
    token: 
    publicKey: 
`)

// ConfYaml is config structure.
type ConfYaml struct {
	Api  api  `yaml:"api"`
	Core core `yaml:"core"`
	Log  log  `yaml:"log"`
}

type core struct {
	DSN string `yaml:"dsn"`
}

// SectionCore is sub section of config.
type api struct {
	Internal internal `yaml:"internal"`
}

type internal struct {
	Token     string `yaml:"token"`
	PublicKey string `yaml:"publicKey"`
}

type log struct {
	Format   string `yaml:"format"`
	ApiLog   string `yaml:"apiLog"`
	ApiLevel string `yaml:"apiLevel"`
}

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath string) (ConfYaml, error) {
	var conf ConfYaml

	viper.SetConfigType("yaml")

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		viper.AddConfigPath("$HOME/go")
		viper.AddConfigPath(".")
		viper.SetConfigName("gochannel_config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	// Core
	conf.Core.DSN = viper.GetString("core.dsn")

	// API
	conf.Api.Internal.Token = viper.GetString("api.internal.token")
	conf.Api.Internal.PublicKey = viper.GetString("api.internal.publicKey")

	// Log
	conf.Log.Format = viper.GetString("log.format")
	conf.Log.ApiLog = viper.GetString("log.apiLog")
	conf.Log.ApiLevel = viper.GetString("log.apiLevel")

	return conf, nil
}
