package main           // a maior parte da lógica da aplicação aqui

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {

	//clientes guarda todos os clientes do room
	clients map[*client]bool

	//join é um canal para o cliente entrar na sala 
	join chan *client

	//leave é um canal para o cliente sair da sala
	leave chan *client

	//forward é um canal que guarda as mensagens recebidas que devem ser enviadas p outros clientes
	forward chan []byte
}

//newRoom cria um novo chatroom
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run(){
	for{
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receptor)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receptor <-msg
			}	
		}

	}
}


const (
	socketBufferSize = 1024
    messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize , WriteBufferSize: messageBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request){
	socket, err := upgrader.Upgrade(w, req,  nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		receptor: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() {r.leave <- client} ()
	go client.write()
	client.read()
} 