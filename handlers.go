package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/strawmang/strawmang/chat"
	"golang.org/x/net/websocket"
)

type Status struct {
	Topics []chat.Topic `json:"topics"`
	Users  int          `json:"users"`
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

// START Fuck CORS

func handlerTest(rw http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadFile("websocket-testing.html")
	if err != nil {
		// TODO: Pretty 503
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(data)
}

func handlerTestJS(rw http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadFile("websocket-testing.js")
	if err != nil {
		// TODO: Pretty 503
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(data)
}

// END Fuck CORS

func handlerWs(rw http.ResponseWriter, req *http.Request) {
	conf, err := websocket.NewConfig("ws://localhost", "http://localhost")
	if err != nil {
		log.Printf("Websocket: %v\n", err.Error())
	}
	s := websocket.Server{Handler: websocket.Handler(chat.HandlerChat), Config: *conf}
	s.ServeHTTP(rw, req)
}

func handlerStatus(rw http.ResponseWriter, req *http.Request) {
	status := new(Status)
	// Collect topics
	status.Topics = []chat.Topic{}
	for _, v := range chat.GlobalServer.Topics {
		status.Topics = append(status.Topics, *v)
	}

	status.Users = len(chat.GlobalServer.Users)

	data, err := json.Marshal(status)
	if err != nil {
		log.Printf("Couldn't marshal JSON: %v\n", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error":"failed to marshal JSON"}`))
	} else {
		log.Println(string(data))
		rw.Write(data)
	}
}

func handlerTestColor(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("templates/colortest.tmpl"))

	colors := []colorful.Color{}

	for i := 0; i < 300; i++ {
		colors = append(colors, chat.PopColor())
	}

	t.Execute(rw, colors)
}
