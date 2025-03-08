package room

import (
	"kotiki/internal/peer"	
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type Room struct {
	mu    sync.Mutex
	peers map[string]*peer.Peer // идентификатор (можно генерировать, как UUID) -> Peer
}

func NewRoom() *Room {
	return &Room{
		peers: make(map[string]*peer.Peer),
	}
}

// Регистрация нового пира
func (r *Room) RegisterPeer(ws *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	peer := peer.NewPeer(ws, r)
	r.peers[peer.ID] = peer
	peer.Listen()

}

// Убираем пира при отключении
func (r *Room) UnregisterPeer(peer *peer.Peer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.peers, peer.ID)
	log.Println("Peer отключился. Осталось участников:", len(r.peers))
}
func (r *Room) BroadcastTrack(track *webrtc.TrackRemote, senderID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, peer := range r.peers {
		if id == senderID {
			continue // не отправлять обратно отправителю
		}

		// Добавляем трек для пересылки
		if err := peer.AddTrack(track); err != nil {
			log.Println("Ошибка добавления трека:", err)
		}
	}
}
// Пересылка track'ов другим участникам (это основа SFU)

