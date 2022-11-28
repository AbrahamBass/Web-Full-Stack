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
	client     []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		client:     make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (this *Hub) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "No se pudo abrir la conexion", http.StatusBadRequest)
	}
	client := NewClient(this, socket)
	this.register <- client

	go client.Write()
}

func (this *Hub) Run() {
	for {
		select {
		case client := <-this.register:
			this.onConnect(client)
		case client := <-this.unregister:
			this.onDisconnect(client)
		}
	}
}

func (this *Hub) onConnect(client *Client) {
	log.Println("Cliente conectado", client.socket.RemoteAddr())

	this.mutex.Lock()
	defer this.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	this.client = append(this.client, client)

}

func (this *Hub) onDisconnect(client *Client) {
	log.Println("Cliente desconectado", client.socket.RemoteAddr())
	client.socket.Close()
	this.mutex.Lock()
	i := -1

	for j, c := range this.client {
		if c.id == client.id {
			i = j
		}
	}
	copy(this.client[i:], this.client[i+1:])
	this.client[len(this.client)-1] = nil
	this.client = this.client[:len(this.client)-1]
}

func (this *Hub) Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range this.client {
		if client != ignore {
			client.outbound <- data
		}
	}
}
