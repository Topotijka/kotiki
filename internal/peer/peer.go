package peer

import (
	
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type Room interface{
		BroadcastTrack(track *webrtc.TrackRemote, senderID string)

}

type Peer struct {
	ID     string
	Conn   *websocket.Conn
	PC     *webrtc.PeerConnection
	signal *Signal
	localTracks map[string]*webrtc.TrackLocalStaticRTP
	room Room
}

type Signal struct {
	Type      string `json:"type"`                // Тип сигнала (offer, answer, candidate)
	SDP       string `json:"sdp,omitempty"`       // SDP предложение или ответ
	Candidate string `json:"candidate,omitempty"` // ICE кандидат
}

func (p *Peer) GetRoom() Room {
	return p.room
}
func NewPeer(conn *websocket.Conn, room Room) *Peer {
	config := webrtc.Configuration{
    ICEServers: []webrtc.ICEServer{
        {
            URLs: []string{"stun:stun.l.google.com:19302"}, // STUN сервер
        },
        // Добавьте TURN сервер, если нужен
        // { urls: 'turn:your-turn-server.com', username: 'user', credential: 'pass' }
    },
}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Println("Ошибка создания PeerConnection:", err)
		return nil
	}
	peer := &Peer{
		ID:         uuid.New().String(),
		Conn:       conn,
		PC:         pc,
		signal:     &Signal{},
		localTracks: make(map[string]*webrtc.TrackLocalStaticRTP),
		room:       room,
	}

	// Обработка входящих треков
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Println("Получен трек:", track.ID(), track.Kind())

		// Пересылка трека другим клиентам
		if room != nil {
			room.BroadcastTrack(track, peer.ID)
		}
	})

	// Отслеживание состояния ICE соединения
	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Printf("Состояние ICE соединения: %s\n", state.String())
	})
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("Состояние соединения: %s\n", state.String())
	})
	return peer
}

func (p *Peer) Listen() {
	defer func() {
		_ = p.PC.Close()
		_ = p.Conn.Close()
	}()

	for {
		err := p.Conn.ReadJSON(&p.signal)
		if err != nil {
			log.Println("Ошибка чтения websocket:", err)
			break
		}

		p.handleSignalMessage()
	}
}

func (p *Peer) handleSignalMessage() {
	switch p.signal.Type {

	case "offer":
		p.handleOffer()
	case "candidate":
		p.handleCandidate()
	case "answer":
		p.handleAnwser()
	}
}

func (p *Peer) forwardTrack(localTrack *webrtc.TrackLocalStaticRTP, remoteTrack *webrtc.TrackRemote) {
	buf := make([]byte, 1500) // Буфер для RTP пакетов

	for {
		// Чтение RTP пакетов из удаленного трека
		n, _, err := remoteTrack.Read(buf)
		if err != nil {
			log.Println("Ошибка чтения RTP пакета:", err)
			return
		}

		// Запись RTP пакетов в локальный трек
		if _, err := localTrack.Write(buf[:n]); err != nil {
			log.Println("Ошибка записи RTP пакета:", err)
			return
		}
	}
}
// Добавление удаленного трека к Peer
func (p *Peer) AddTrack(track *webrtc.TrackRemote) error {
	localTrack, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
	if err != nil {
		return err
	}

	// Сохраняем локальный трек
	p.localTracks[track.ID()] = localTrack

	// Добавляем локальный трек в PeerConnection
	_, err = p.PC.AddTrack(localTrack)
	if err != nil {
		return err
	}

	// Запускаем пересылку данных
	go p.forwardTrack(localTrack, track)

	return nil
}
