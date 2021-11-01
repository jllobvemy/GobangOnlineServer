package main

import (
	"log"
)

type SingleClient struct {
	id string
	client *Client
}

type Message struct{
	id string
	message []byte
}

type Hub struct {
	clients map[string] []*Client
	broadcast chan Message
	register chan SingleClient
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan Message),
		register: make(chan SingleClient),
		unregister: make(chan *Client),
		clients: make(map[string] []*Client),
	}
}

func (h *Hub) run()  {
	for {
		select {
			case client := <-h.register:
				h.clients[client.id] = append(h.clients[client.id], client.client)
				log.Println("Connected: " + client.id)
			case client := <-h.unregister:
				for k, currClient := range h.clients{
					for i := 0; i < len(currClient); i++ {
						if client == currClient[i]{
							h.clients[k] = append(currClient[:i], currClient[i+1:]...)
							log.Println("Disconnected: " + client.id)
						}
					}
				}
			case message := <-h.broadcast:
				for _, client := range h.clients[message.id]{
					client.send <- message.message
				}
		}
	}
}
