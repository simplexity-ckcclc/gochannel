package main

import (
	"github.com/ckcclc-simplexity/gochannel/api"
	"github.com/ckcclc-simplexity/gochannel/match"
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
