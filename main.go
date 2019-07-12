package main

import (
	"flag"
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match"
)

func main() {
	//tools.SignAndPrintNonce()
	//tools.SignAndPrintClickInfo()

	var confPath string
	flag.StringVar(&confPath, "conf", "", "gochannel config file")
	flag.Parse()

	// load config
	var err error
	if err = common.LoadConf(confPath); err != nil {
		panic(err)
	}

	// start match-Server
	go match.Serve()

	// start api-server
	api.Serve()

}
