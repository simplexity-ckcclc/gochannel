package main

import (
	"github.com/simplexity-ckcclc/gochannel/api"
	"github.com/simplexity-ckcclc/gochannel/match"
	"net/http"
)

func main() {
	// start Match-Server
	go match.Serve()

	server := &http.Server{
		Addr:    ":8480",
		Handler: api.Router(),
	}
	server.ListenAndServe()

}
