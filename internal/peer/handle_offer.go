package peer

import (
	"log"

	"github.com/pion/webrtc/v4"
)

func (p *Peer) handleOffer() {

	offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: p.signal.SDP}
	err := p.PC.SetRemoteDescription(offer)
	if err != nil {
		log.Println("Ошибка установки remote description:", err)
		return
	}

	// Создать ответ (answer) и отправить клиенту
	answer, err := p.PC.CreateAnswer(nil)
	if err != nil {
		log.Println("Ошибка создания ответа:", err)
		return

	}

	err = p.PC.SetLocalDescription(answer)
	if err != nil {
		log.Println("Ошибка установки local description:", err)
		return

	}

	answerMessage := Signal{
		Type: "answer",
		SDP:  answer.SDP,
	}
	log.Println(answerMessage, "/n")
	p.Conn.WriteJSON(answerMessage)

}
