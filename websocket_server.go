package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	server := NewServer()
	http.HandleFunc("/", handler(server))
	log.Println("Server starting on :8008")
	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		panic(err)
	}
}

type Client struct {
	conn *websocket.Conn
}

type Server struct {
	clients map[Client]struct{}

	addClient chan Client
	removeClient chan Client
	messages chan Message
}

type Message struct {
	client Client
	message string
}

func NewServer() *Server {
	server := &Server{
		clients: make(map[Client]struct{}),
		addClient: make(chan Client),
		removeClient: make(chan Client),
		messages: make(chan Message),
	}
	go server.Start()
	return server
}

func (s *Server) Start() {
	for {
		select {
			case client := <-s.addClient:
				s.clients[client] = struct{}{}
			case client := <-s.removeClient:
				delete(s.clients, client)
			case message := <-s.messages:
				log.Println("<-", message.message)
				for client := range s.clients {
					if client == message.client {
						continue
					}
					err := client.conn.WriteMessage(websocket.TextMessage, []byte(message.message))
					if err != nil {
						log.Println("error writing:", err)
					}
				}
		}
	}
}


func handler(server *Server) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  10,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := Client{
			conn,
		}
		server.addClient <- client

		for {
			messageType, r, err := conn.NextReader()
			if err != nil {
				log.Println(err)
				server.removeClient <- client
				return
			}
			if messageType == websocket.TextMessage {
				bytes, err := ioutil.ReadAll(r)
				if err != nil {
					log.Println(err)
					server.removeClient <- client
					return
				}
				message := string(bytes)
				server.messages <- Message{
					client:  client,
					message: message,
				}
			}
		}
	}
}
