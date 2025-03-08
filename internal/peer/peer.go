package peer

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type Peer struct {
	ID     string
	Conn   *websocket.Conn
	PC     *webrtc.PeerConnection
	signal *Signal
}

type Signal struct {
	Type      string `json:"type"`                // Тип сигнала (offer, answer, candidate)
	SDP       string `json:"sdp,omitempty"`       // SDP предложение или ответ
	Candidate string `json:"candidate,omitempty"` // ICE кандидат
}

func NewPeer(conn *websocket.Conn) *Peer {

	config := webrtc.Configuration{}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Println("cant create peer conn")
	}
	peer := &Peer{
		ID:     uuid.New().String(),
		Conn:   conn,
		PC:     pc,
		signal: &Signal{},
	}

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

// Добавление удаленного трека к Peer
func (p *Peer) AddTrack(track *webrtc.TrackRemote) error {
	localTrack, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
	if err != nil {
		return err
	}
	_, err = p.PC.AddTrack(localTrack)
	return err
}
