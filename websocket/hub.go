package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "You Could not open websocket connection", http.StatusBadRequest)
	}
	client := NewClient(hub, socket)
	hub.register <- client
	go client.Write()
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnected(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnected(client *Client) {
	log.Println("Client Connected", client.socket.RemoteAddr())

	hub.mutex.Lock() //Lock the programs for no ruuner conditions
	defer hub.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client) //Add new client into existing clients
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client Disconnected", client.socket.RemoteAddr())
	client.socket.Close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
		}
	}
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]
}

func (hub *Hub) BroadCast(message interface{}, ingnore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ingnore {
			client.outbound <- data
		}
	}
}
