package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"time"
)

func handlerChat(conn *websocket.Conn) {
	// reuse buffers;  keep memory usage low!
	var buff [1024]byte
	var event Event

	var me *User
	// TODO: This loop is going to get very unweildy.  Break it up
loop:
	for {
		n, err := conn.Read(&buff)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Printf("[chat] EOF\n")
				conn.Close()
				break loop
			}

			server.Errors <- err
		}

		if err := json.Unmarshal(buff[:n], &event); err != nil {
			server.Errors <- err
		}
		switch event.Type {
		case EVENT_MESSAGE:
			if me != nil {
				server.Events <- event
			}
		case EVENT_LEAVE:
			server.Leave(me)
		case EVENT_LOGIN:
			if me != nil {
				SendMessage(conn, NewErrorEvent(errors.New("Can't login twice, dumbass")))
			} else {
				me = new(User)
				me.conn = conn
				me.Name = event.Username
				me.Remote = conn.RemoteAddr().String()
				me.Color = generateColor()
				if err := server.Join(me); err != nil {
					SendMessage(conn, NewErrorEvent(err))
					conn.Close()
					break loop
				}
				me.ListenEvents()
				log.Printf("Login worked.  Listening for events")
			}
		}
	}
}

// TODO: This
func generateColor() string {
	return "000000"
}

// getNextTopicID returns the next topic ID
func getNextTopicID() int {
	i := 1
	// do logic
	return i
}

// Event is a set of data that will be sent over the websocket
//
// Types:
// login:
// message:
// leave:
// vote:
//
// TODO: Set omitempty on sane choices
type Event struct {
	// global
	Type string `json:"type"`

	// login
	Username string `json:"username"`

	// message
	Text    string `json:"text"`
	TopicID string `json:"topic-id"`

	source string `json:"-"` // The remote address of the user
	Error  string `json:"error"`
}

func SendEvent(conn *websocket.Conn, ev Event) {
	data, err := json.Marshal(ev)
	if err != nil {
		server <- err
	}

	if _, err := conn.Write(data); err != nil {
		server <- err
	}
}

func NewErrorEvent(err error) Event {
	return Event{Error: err}
}

var server = new(Server)

// Server represents a single chat server that will run
type Server struct {
	Events chan Event
	Errors chan error
	Kill   chan struct{}
	Users  map[string]*User
}

// Join should be called when a user wants to join a server
func (s *Server) Join(user *User) error {
	if user == nil {
		return errors.New("server: Join called will nil user? This is a bug")
	}

	// @TODO This probably needs to be mutexted?
	s.Users[user.conn.RemoteAddr().String()] = user

	return nil
}

// Leave deletes the user instance from the server if it exists
func (s *Server) Leave(user *User) {
	if user == nil {
		return
	}

	if _, ok := s.Users[user.conn.RemoteAddr().String()]; ok {
		delete(s.Users, user.conn.RemoteAddr().String())
	}
}

// Start will start all of the event passing logic
func (s *Server) Start() error {
	go func() {
		for {
			select {
			case err := <-s.Errors:
				log.Printf("server error: %v", err.Error())
			case <-s.Kill:
				log.Printf("server logic shutting down")
			case event := <-s.Events:
				for _, v := range s.Users {
					v.Events <- event
				}
			}
		}
	}()
}

// Stop will attempt to gracely stop the server and close the topic
func (s *Server) Stop() {
	s.Kill <- struct{}{}
}

// User is a single user that can be connected to a user
// When a user joins a channel a goroutine is started to listen for
// data on the websocket
type User struct {
	Name   string          `json:"name"`
	Color  string          `json:"color"`
	Remote string          `json:"-"`
	conn   *websocket.Conn `json:"-"`
	Events chan Event      `json:"-"`
	Kill   chan struct{}   `json:"-"`
}

func (u *User) ListenEvents() {
	go func() {
		var event Event
	loop:
		for {
			select {
			case event = <-u.Events:
				data, err := json.Marshal(event)
				if err != nil {
					server.Errors <- err
				}
				u.conn.Write(data) // Ignore errors.
			case <-u.Kill:
				break loop
			}
		}
	}()
}

// Topic is a single topic that can be ran on a server
type Topic struct {
	OptionA string `json:"option-a"`
	OptionB string `json:"option-b"`

	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	Ends    time.Time `json:"ends"`

	events chan Event `json:"-"`
}
