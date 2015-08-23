package main

import (
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

// getNextTopicID returns the next topic ID
func getNextTopicID() int {
	i := 1
	return i
}

// Event is a set of data that will be sent over the websocket
//
// Types:
// message:
// leave:
// vote:
type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Server represents a single chat server that will run
type Server struct {
}

// Join should be called when a user wants to join a server
func (s *Server) Join(user *User) error {}

// Start will allow the server to start listening to connections
func (s *Server) Start() error {}

// Stop will attempt to gracely stop the server and close the topic
func (s *Server) Stop() error {}

// User is a single user that can be connected to a user
// When a user joins a channel a goroutine is started to listen for
// data on the websocket
type User struct {
	Name  string   `json:"name"`
	IP    string   `json:"color"`
	Color string   `json:"color"`
	conn  *ws.Conn `json:"-"`
}

// Listen starts a goroutine to listen for data on the websocket
func (u *User) Listen() {}

// Topic is a single topic that can be ran on a server
type Topic struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	Ends    time.Time `json:"ends"`

	events chan Event `json:"-"`
}
