package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

// Room represents a chat room.
type Room struct {
	Name     string
	Clients  []*melody.Session
	Messages []string
}

// Rooms stores all the chat rooms.
var Rooms = make(map[string]*Room)

func containsClients(Clients []*melody.Session, target *melody.Session) bool {
	for _, num := range Clients {
		if num == target {
			return true
		}
	}
	return false
}

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws/:room", func(c *gin.Context) {
		roomName := c.Param("room")
		if _, ok := Rooms[roomName]; !ok {
			Rooms[roomName] = &Room{Name: roomName}
			fmt.Println("room added", Rooms)
		} else {
			fmt.Println("can't find room")
		}
		m.HandleRequest(c.Writer, c.Request)

	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		path := s.Request.URL.Path
		path = strings.Replace(path, "/ws/", "", 1)
		fmt.Println(path)
		fmt.Println(Rooms[path])
		//add client sessions to room

		if !containsClients(Rooms[path].Clients, s) {

			Rooms[path].Clients = append(Rooms[path].Clients, s)
		}

		for _, client := range Rooms[path].Clients {
			client.Write(msg)
		}

		//add message to room
		Rooms[path].Messages = append(Rooms[path].Messages, string(msg))
	})

	r.Run(":5000")
}
