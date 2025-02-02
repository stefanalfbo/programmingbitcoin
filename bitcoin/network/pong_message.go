package network

import (
	"encoding/binary"
	"io"
)

// The pong message is sent in response to a ping message. In modern protocol
// versions, a pong response is generated using a nonce included in the ping.
type PongMessage struct {
	command []byte
	// Nonce from ping
	Nonce uint64
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

func (pm *PongMessage) Parse(reader io.Reader) (Message, error) {
	var nonce uint64
	err := binary.Read(reader, binary.LittleEndian, &nonce)
	if err != nil {
		return nil, err
	}

	return NewPongMessage(nonce), nil
}
