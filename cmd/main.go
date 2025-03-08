package main

import (
	"kotiki/internal/room"
	"kotiki/internal/server"
	"log"
	"net/http"
)

func main() {
	room := room.NewRoom()           // создаем единственную комнату
	server := server.NewServer(room) // создаем сервер

	http.HandleFunc("/ws", server.HandleWebSocket)

	log.Println("SFU сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
