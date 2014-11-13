package main

import (
	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	add = make(chan *Client)
)

type Room struct {
	Clients map[string]*Client
}

func (r *Room) Add(id string, c *Client) {
	r.Clients[id] = c
}

func Broadcast(excludeID string, clients map[string]*Client) {
	for id, c := range clients {
		if id == excludeID {
			continue
		}

		websocket.JSON.Send(c.Ws, Action{
			Type: "otherCreate",
			ID:   id,
			X:    c.Action.X,
			Y:    c.Action.Y,
		})
	}
}

func (r *Room) Broadcast() {
	for uid, _ := range r.Clients {
		go Broadcast(uid, r.Clients)
	}
}

type Action struct {
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type Client struct {
	Action Action
	Ws     *websocket.Conn
}

func openRoom() {
	var room = &Room{
		make(map[string]*Client),
	}
	for {
		select {
		case client := <-add:
			id := uuid.NewUUID().String()
			room.Add(id, client)
			room.Broadcast()

			// generate id to client
			websocket.JSON.Send(client.Ws, map[string]string{
				"type": "create",
				"id":   id,
			})
		}
		log.Println(room)
	}
}

func dispatch(ws *websocket.Conn, action Action) {
	switch action.Type {
	case "create":
		add <- &Client{action, ws}
	}
}

func RoomHandler(ws *websocket.Conn) {
	var action Action
	for websocket.JSON.Receive(ws, &action) != io.EOF {
		go dispatch(ws, action)
	}
}

func verbose(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func main() {
	var port string
	if len(os.Args) != 2 {
		port = "8100"
	} else {
		port = os.Args[1]
	}

	go openRoom()

	http.Handle("/", verbose(http.FileServer(http.Dir("www"))))
	http.Handle("/join", websocket.Handler(RoomHandler))

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}
