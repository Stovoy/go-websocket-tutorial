package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8008", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.TextMessage {
			bytes, err := ioutil.ReadAll(r)
			if err != nil {
				log.Println(err)
				return
			}
			text := string(bytes)
			log.Println("<-", text)
		}
	}
}
