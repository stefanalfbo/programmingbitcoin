package network

import (
	"encoding/binary"
	"io"
)

// The ping message is sent primarily to confirm that the TCP/IP connection
// is still valid. An error in transmission is presumed to be a closed
// connection and the address is removed as a current peer.
type PingMessage struct {
	command []byte
	// Random nonce
	Nonce uint64
}

func NewPingMessage(nonce uint64) *PingMessage {
	return &PingMessage{
		command: []byte("ping"),
		Nonce:   nonce,
	}
}

func (pm *PingMessage) Command() []byte {
	return pm.command
}

func (pm *PingMessage) Serialize() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, pm.Nonce)
	return buf, nil
}

func (pm *PingMessage) Parse(reader io.Reader) (Message, error) {
	var nonce uint64
	err := binary.Read(reader, binary.LittleEndian, &nonce)
	if err != nil {
		return nil, err
	}

	return NewPingMessage(nonce), nil
}
