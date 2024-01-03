package main

import (
	"html/template"
	"net/http"

	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var conns = map[string]*websocket.Conn{}

type Message struct {
	Text string `json:"Text"`
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(conns, conn.RemoteAddr().String())
			conn.Close()
			return
		}

		fmt.Printf("message from client: %s\n", string(p))

		msg := fmt.Sprintf("%v wrote: %s", conn.RemoteAddr(), string(p))

		for _, v := range conns {
			if err := v.WriteMessage(messageType, []byte(msg)); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		messages := map[string][]Message{
			"Messages": {
				{Text: "first message"},
			},
		}
		tmpl.Execute(w, messages)
	}

	ws := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("conn received")
		upgrade.CheckOrigin = func(r *http.Request) bool { return true }

		wsc, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer wsc.Close()

		conns[wsc.RemoteAddr().String()] = wsc

		reader(wsc)
	}

	http.HandleFunc("/ws", ws)
	http.HandleFunc("/home", h1)
	http.ListenAndServe(":8080", nil)

}
