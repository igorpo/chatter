package chatroom

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type Room struct {
	messageQueue chan []byte
	joinQueue    chan *client // allows for safe modification of our client set for joining
	leaveQueue   chan *client // and leaving the room
	clients      map[*client]bool
}

func NewRoom() *Room {
	return &Room{
		messageQueue: make(chan []byte),
		joinQueue:    make(chan *client),
		leaveQueue:   make(chan *client),
		clients:      make(map[*client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.joinQueue:
			r.clients[client] = true
		case client := <-r.leaveQueue:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.messageQueue:
			// broadcast to all clients
			for c := range r.clients {
				c.send <- msg
			}
		}
	}
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.joinQueue <- client
	defer func() { r.leaveQueue <- client }()
	go client.write()
	client.read()
}
