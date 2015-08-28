package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/strawmang/strawmang/chat"
)

// Dev is true when we detect a development environment the default
// stratagy is to check for .git in the current directory
var Dev bool

// Commit is the current commit strawmang was built on
var Commit string

func init() {
	if _, err := os.Stat(".git"); err == nil {
		Dev = true
		log.Printf("strawmang server running in Dev mode")
	}
	if _, err := os.Stat("index.html"); err != nil {
		log.Panic("No index.html found;  fix the deployment")
	}
}

func main() {
	chat.GlobalServer.Start()
	r := mux.NewRouter()

	r.HandleFunc("/", handlerIndex)
	r.HandleFunc("/ws", handlerWs)
	r.HandleFunc("/status", handlerStatus)

	// Maybe: https://groups.google.com/forum/#!topic/golang-nuts/bStLPdIVM6w ?
	//      : Removing directory listing

	// TODO: Route from memory using go-bindata
	//     : https://github.com/elazarl/go-bindata-assetfs

	fs := http.FileServer(http.Dir("static/"))

	r.PathPrefix("/").Handler(fs)

	if Dev {
		r.HandleFunc("/test", handlerTest)
		r.HandleFunc("/websocket-testing.js", handlerTestJS)

		r.HandleFunc("/colortest", handlerTestColor)
	}

	if Commit != "" {
		log.Printf("  Running commit %s", Commit)
	}

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("http: %v", err.Error())
	}
}
