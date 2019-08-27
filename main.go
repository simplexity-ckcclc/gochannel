package main

import (
	"flag"
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
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
	if err = config.LoadConf(confPath); err != nil {
		panic(err)
	}

	if err = initiate(); err != nil {
		panic(err)
	}
	defer destroy()

	// start match-Server
	go match.Serve()

	// start api-server
	go api.Serve()

	running := make(chan bool, 1)
	<-running
}

func initiate() (err error) {
	if err = common.InitLogger(); err != nil {
		return
	}

	if _, err = common.InitSqlClient(); err != nil {
		return
	}

	// init es client
	if _, err = common.InitEsClient(); err != nil {
		return
	}

	return
}

func destroy() {
	_ = common.DB.Close()
	common.EsClient.Stop()
}
