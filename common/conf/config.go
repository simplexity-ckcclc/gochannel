package config

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
	Api  API  `yaml:"api"`
	Core Core `yaml:"core"`
}

type Core struct {
	DSN string `yaml:"dsn"`
}

// SectionCore is sub section of config.
type API struct {
	Internal Internal `yaml:"internal"`
}

type Internal struct {
	Token     string `yaml:"token"`
	PublicKey string `yaml:"publicKey"`
}

type Common struct {
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

	return conf, nil
}
