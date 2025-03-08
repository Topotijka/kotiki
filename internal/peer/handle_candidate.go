package peer

import (
	"fmt"
	"log"

	"github.com/pion/webrtc/v4"
)

func (p *Peer) handleCandidate() {
	fmt.Print("''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''''")
	candidate := webrtc.ICECandidateInit{Candidate: p.signal.Candidate}
	err := p.PC.AddICECandidate(candidate)
	if err != nil {
		log.Println("candidate error ", err.Error())
	}

}
