package main

import (
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/api/entity"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match"
	"net/http"
)

const DSN = "ckcclc:141421@tcp(localhost:3306)/gochannel"

func main() {
	// start Match-Server
	go match.Serve()

	//db, err := sql.Open("mysql", DSN)
	db, err := common.OpenDB(DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := entity.LoadAppKeySigs(db); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    ":8480",
		Handler: api.Router(),
	}
	server.ListenAndServe()

}
