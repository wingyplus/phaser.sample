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
	add  = make(chan *Client)
	move = make(chan MoveAction)
)

const (
	OTHER_CREATE = "otherCreate"
	MOVE         = "move"
	CREATE       = "create"
)

type Room struct {
	Clients map[string]*Client
}

func (r *Room) Add(id string, c *Client) {
	r.Clients[id] = c
}

func Broadcast(excludeID string, clients map[string]*Client, event string) {
	for id, c := range clients {
		if id == excludeID {
			continue
		}

		switch event {
		case OTHER_CREATE:
			println("do otherCreate")
			websocket.JSON.Send(c.Ws, Action{
				Type: OTHER_CREATE,
				ID:   id,
				X:    c.Position.X,
				Y:    c.Position.Y,
			})
		case MOVE:
			websocket.JSON.Send(c.Ws, Action{
				Type: MOVE,
				ID:   id,
				X:    c.Velocity.X,
				Y:    c.Velocity.Y,
			})
		}
	}
}

func (r *Room) Broadcast(event string) {
	for uid, _ := range r.Clients {
		go Broadcast(uid, r.Clients, event)
	}
}

type Action struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type MoveAction struct {
	ID       string
	Velocity Velocity
}

type Position struct {
	X, Y int
}

type Velocity struct {
	X, Y int
}

type Client struct {
	Position Position
	Velocity Velocity
	Ws       *websocket.Conn
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
			go room.Broadcast(OTHER_CREATE)

			// generate id to client
			websocket.JSON.Send(client.Ws, map[string]string{
				"type": CREATE,
				"id":   id,
			})
		case action := <-move:
			room.Clients[action.ID].Velocity = action.Velocity
			room.Broadcast(MOVE)
		}
		log.Println(room)
	}
}

func dispatch(ws *websocket.Conn, action Action) {
	println(action.Type)
	switch action.Type {
	case CREATE:
		add <- &Client{Position: Position{action.X, action.Y}, Ws: ws}
	case MOVE:
		move <- MoveAction{ID: action.ID, Velocity: Velocity{action.X, action.Y}}
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
