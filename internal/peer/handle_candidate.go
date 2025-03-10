package peer

import (

	"log"

	"github.com/pion/webrtc/v4"
)

func (p *Peer) handleCandidate() {
    candidate := webrtc.ICECandidateInit{Candidate: p.signal.Candidate}
    err := p.PC.AddICECandidate(candidate)
    if err != nil {
        log.Println("Ошибка добавления ICE кандидата:", err)
    }
}
