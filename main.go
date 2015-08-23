package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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
	r := mux.NewRouter()

	r.HandleFunc("/", handlerIndex)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("http: %v", err.Error())
	}
}
