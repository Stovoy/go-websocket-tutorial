package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

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

	go func() {
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
	}()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		message := sc.Text()
		log.Println("->", message)
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			panic(err)
		}
	}

	if sc.Err() != nil {
		panic(err)
	}
}
