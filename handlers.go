package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/strawmang/strawmang/chat"
	"golang.org/x/net/websocket"
)

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
func handlerWs(rw http.ResponseWriter, req *http.Request) {
	conf, err := websocket.NewConfig("ws://localhost", "http://localhost")
	if err != nil {
		log.Printf("Websocket: %v\n", err.Error())
	}
	s := websocket.Server{Handler: websocket.Handler(chat.HandlerChat), Config: *conf}
	s.ServeHTTP(rw, req)
}
