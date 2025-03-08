package server

import (
	"kotiki/internal/room"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	Room     *room.Room
	upgrader websocket.Upgrader
}

func NewServer(room *room.Room) *Server {
	return &Server{
		Room: room,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при подключении websocket:", err)
		return
	}
	defer wsConn.Close()

	s.Room.RegisterPeer(wsConn)

}
