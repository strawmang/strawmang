package chat

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"sync"
	"time"
)

// TODO: The server needs a better way to identify a person
//       based off of their connection.  Currently if multiple
//       people login from the same IP there will be strange
//       things occuring

const MaxTopics = 3

var GlobalServer = NewServer()

// Event is a set of data that will be sent over the websocket
//
// See events.md for more information
type Event struct {
	// global
	Type string `json:"type"`

	// login
	Username string `json:"username,omitempty"`
	Color    string `json:"color,omitempty"`

	// message
	Text    string `json:"text,omitempty"`
	TopicID string `json:"topic-id,omitempty"`

	// newtopic
	OptionA string `json:"option-a"`
	OptionB string `json:"option-b"`

	source string `json:"-"` // The remote address of the user
	Error  string `json:"error,omitempty"`
}

// Server represents a single chat server that will run
type Server struct {
	Users  map[string]*User `json:"-"`
	Topics map[int]*Topic   `json:"users"`

	// internal shit
	topicsMutext *sync.Mutex

	Errors chan error    `json:"-"`
	Kill   chan struct{} `json:"-"`
	Events chan Event    `json:"-"`
}

// User is a single user that can be connected to a user
// When a user joins a channel a goroutine is started to listen for
// data on the websocket
type User struct {
	Name     string          `json:"name"`
	Color    string          `json:"color"`
	Remote   string          `json:"-"`
	conn     *websocket.Conn `json:"-"`
	Events   chan Event      `json:"-"`
	Kill     chan struct{}   `json:"-"`
	LoggedIn bool            `json:"logged-in"`
}

// Topic is a single topic that can be ran on a server
type Topic struct {
	OptionA string `json:"option-a"`
	OptionB string `json:"option-b"`

	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	Ends    time.Time `json:"ends"`

	events chan Event `json:"-"`
}

// Join should be called when a user wants to join a server
func (s *Server) Join(user *User) error {
	if user == nil {
		return errors.New("GlobalServer: Join called will nil user? This is a bug")
	}

	log.Printf("=> %v successfully joined!", user.Name)
	// @TODO This probably needs to be mutexted?
	s.Users[user.conn.RemoteAddr().String()] = user
	user.Color = generateColor()

	return nil
}

// Leave deletes the user instance from the server if it exists
func (s *Server) Leave(user *User) {
	if user == nil {
		return
	}

	if _, ok := s.Users[user.conn.RemoteAddr().String()]; ok {
		SendEvent(s.Users[user.conn.RemoteAddr().String()].conn, Event{Type: EVENT_STATUS, Text: "kbye"})

		log.Printf("<= %v successfully left!", user.Name)
		delete(s.Users, user.conn.RemoteAddr().String())
		user.Kill <- struct{}{}
		//user.conn.Close()
	}
}

// Start will start all of the event passing logic
func (s *Server) Start() error {
	go func() {
	loop:
		for {
			select {
			case err := <-s.Errors:
				log.Printf("Server error: %v", err.Error())
			case <-s.Kill:
				log.Printf("Server logic shutting down")
				break loop
			case event := <-s.Events:
				// Normally only send if topic ID is active or ID is -1 for testing purposes
				for _, v := range s.Users {
					v.Events <- event
				}
			default:
			}
		}
	}()
	return nil
}

// Stop will attempt to gracely stop the server and close the topic
func (s *Server) Stop() {
	s.Kill <- struct{}{}
}

// TryNewTopic attempts to create a new topic if we have less than three
// topics.
//
// This method takes advantage of a sync.Mutex
func (s *Server) TryNewTopic(t Topic) error {
	s.topicsMutext.Lock()
	defer s.topicsMutext.Unlock()
	if len(s.Topics) >= MaxTopics {
		return errors.New("too many topics already")
	}
	id := s.newTopicID()
	s.Topics[id] = &t
	return nil
}

func (s *Server) newTopicID() int {
	i := 1
	for k, _ := range s.Topics {
		if k > i {
			i = k + 1
		}
	}
	return i
}

// TODO: Make all the handlers send an error if not logged in

func (u *User) HandleNewTopic(event Event) {
	if u.LoggedIn {
		// TODO: Check for whitespace and other weird characters.  Don't want any blank arguments.
		if event.OptionA == "" || event.OptionB == "" {
			SendEvent(u.conn, NewErrorEvent(errors.New("Can't have a blank topic!")))
			return
		}
		t := new(Topic)
		t.OptionA = event.OptionA
		t.OptionB = event.OptionB
		if err := GlobalServer.TryNewTopic(*t); err != nil {
			SendEvent(u.conn, NewErrorEvent(err))
			return
		}
		SendEvent(u.conn, Event{Type: EVENT_STATUS, Text: "Successfully added topic"})
	}
}

func (u *User) HandleLogin(event Event) {
	if u.LoggedIn {
		SendEvent(u.conn, NewErrorEvent(errors.New("Can't login twice, dumbass")))
	} else {
		u.LoggedIn = true
		u.Name = event.Username
		u.Color = generateColor()
		u.Events = make(chan Event)
		if err := GlobalServer.Join(u); err != nil {
			SendEvent(u.conn, NewErrorEvent(err))
		}
		u.ListenEvents()
		SendEvent(u.conn, Event{Type: EVENT_STATUS, Text: "Login ok"})
	}
}
func (u *User) HandleLeave(event Event) {
	GlobalServer.Leave(u)
}
func (u *User) HandleMessage(event Event) {
	// TODO: Send an error if the user is not logged in
	if u.LoggedIn {
		event.Color = u.Color
		GlobalServer.Events <- event
	}
}

// ListenEvents starts a goroutine to read for events.  If stops
// if you send anything to User.Kill
func (u *User) ListenEvents() {
	go func() {
		var event Event
	loop:
		for {
			select {
			case event = <-u.Events:
				data, err := json.Marshal(event)
				if err != nil {
					GlobalServer.Errors <- err
				}
				u.conn.Write(data) // Ignore errors.
			case <-u.Kill:
				break loop
			default:
			}
		}
	}()
}

func HandlerChat(conn *websocket.Conn) {
	defer log.Printf("Closing websocket connection")
	log.Printf("New websocket connection")

	// reuse buffers;  keep memory usage low!
	buff := make([]byte, 1024)
	var event Event
	me := new(User)
	me.conn = conn
	me.Remote = conn.RemoteAddr().String()

loop:
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Printf("[chat] EOF\n")
				conn.Close()
				break loop
			}
			GlobalServer.Errors <- err
		}

		log.Printf("Got: '%v'", string(buff[:n]))

		if err := json.Unmarshal(buff[:n], &event); err != nil {
			GlobalServer.Errors <- err
		}
		switch event.Type {
		case EVENT_NEWTOPIC:
			me.HandleNewTopic(event)
		case EVENT_MESSAGE:
			me.HandleMessage(event)
		case EVENT_LEAVE:
			me.HandleLeave(event)
		case EVENT_LOGIN:
			me.HandleLogin(event)
		default:
			log.Printf("Unhanlded event type in user handler")
		}
	}
}

// NewErrorEvent returns an Event with the error field filled for you
func NewErrorEvent(err error) Event {
	return Event{Type: EVENT_ERROR, Error: err.Error()}
}

// NewServer creates a new Server with all of the fieldsd initialized
func NewServer() *Server {
	s := new(Server)
	s.Events = make(chan Event, 20)
	s.Errors = make(chan error, 5)
	s.Kill = make(chan struct{})
	s.Users = map[string]*User{}

	s.topicsMutext = new(sync.Mutex)
	return s
}

// SendEvent will marshal an event to JSON and write it to the connection
// it passes all errors to the Server error channel
func SendEvent(conn *websocket.Conn, ev Event) {
	data, err := json.Marshal(ev)
	if err != nil {
		GlobalServer.Errors <- err
	}

	if _, err := conn.Write(data); err != nil {
		GlobalServer.Errors <- err
	}
}

func SendEventClient(conn *websocket.Conn, ev Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	if _, err := conn.Write(data); err != nil {
		return err
	}
	return nil
}
