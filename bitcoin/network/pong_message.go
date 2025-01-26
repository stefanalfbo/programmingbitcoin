package network

import (
	"encoding/binary"
	"fmt"
)

type PongMessage struct {
	command []byte
	Nonce   uint64
}

func NewPongMessage(nonce uint64) *PongMessage {
	return &PongMessage{
		command: []byte("pong"),
		Nonce:   nonce,
	}
}

func (pm *PongMessage) Command() []byte {
	return pm.command
}

func (pm *PongMessage) Serialize() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, pm.Nonce)
	return buf, nil
}

func (pm *PongMessage) Parse(data []byte) (Message, error) {
	if len(data) != 8 {
		return nil, fmt.Errorf("invalid pong message length")
	}

	nonce := binary.LittleEndian.Uint64(data)
	return NewPongMessage(nonce), nil
}
