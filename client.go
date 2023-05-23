package main

import (
	"github.com/gorilla/websocket"
	
)
//Cliente representa a conexão com o nosso user

type client struct{

	


	// socket é o web socket para este cliente representa a conexão que utilizamos para mandar e receber mensagens do cliente
	socket *websocket.Conn

	//receptor é um canal pra receber mensagens de outros clientes
	receptor chan []byte 
	//room é o o room onde este cliente está usando o chat
	room *room 

}

// métodos do cliente

func (c *client) read(){
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil{
			return
		}
		c.room.forward <- msg

	}
}

func(c *client) write(){
	defer c.socket.Close()
	for msg := range c.receptor {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}


