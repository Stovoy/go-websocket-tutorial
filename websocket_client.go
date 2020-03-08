package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	client := websocket.Dialer{}
	conn, resp, err :=  client.Dial("ws://localhost:8008", nil)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if len(bytes) > 0 {
		log.Println("Body: ", string(bytes))
	}
	for {
		message := "message"
		log.Println("->", message)
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
