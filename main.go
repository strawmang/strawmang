package main

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Dev is true when we detect a development environment the default
// stratagy is to check for .git in the current directory
var Dev bool

func init() {
	if _, err := os.Stat(".git"); err == nil {
		Dev = true
		log.Printf("strawmang server running in Dev mode")
	}
	if _, err := os.Stat("index.html"); err != nil {
		log.Panic("No index.html found;  fix the deployment")
	}
}

// TODO: Not production ready.  Needs to save the index in memory and only reload it if the file changes
func handlerIndex(rw http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadFile("index.html")
	if err != nil {
		// TODO: Pretty 503
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(data)
}

func main() {
	server.Start()
	r := mux.NewRouter()

	r.HandleFunc("/", handlerIndex)
	r.HandleFunc("/ws", func(rw http.ResponseWriter, req *http.Request) {
		conf, err := websocket.NewConfig("ws://localhost", "http://localhost")
		if err != nil {
			log.Printf("Websocket: %v\n", err.Error())
		}
		s := websocket.Server{Handler: websocket.Handler(handlerChat), Config: *conf}
		s.ServeHTTP(rw, req)
	})
	//r.Handle("/ws", websocket.Handler(handlerChat))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("http: %v", err.Error())
	}
}
