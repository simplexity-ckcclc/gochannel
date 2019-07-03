package main

import (
	"flag"
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match"
	"net/http"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "gochannel config file")

	flag.Parse()

	// load config
	var err error
	common.Conf, err = common.LoadConf(confPath)
	if err != nil {
		panic(err)
	}

	if err := common.InitLog(common.Conf); err != nil {
		panic(err)
	}

	db, err := common.OpenDB(common.Conf.Core.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// start match-Server
	go match.Serve()

	// start api-server
	if err := entity.LoadAppKeySigs(db); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    ":8480",
		Handler: api.Router(),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
