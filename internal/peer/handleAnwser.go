package peer

import (
	"log"

	"github.com/pion/webrtc/v4"
)

func (p *Peer) handleAnwser() {
	log.Println("ans")

	answer := webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: p.signal.SDP}
	err := p.PC.SetRemoteDescription(answer)
	if err != nil {
		log.Println("Ошибка установки remote description:", err)
	}
}
