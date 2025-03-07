package main

import (
	"kotiki/internal"
	"log"
	"net/http"
)

func main() {
	room := internal.NewRoom()         // создаем единственную комнату
	server := internal.NewServer(room) // создаем сервер

	http.HandleFunc("/ws", server.HandleWebSocket)

	log.Println("SFU сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
