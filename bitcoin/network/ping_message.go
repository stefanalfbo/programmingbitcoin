package network

import "encoding/binary"

type PingMessage struct {
	command []byte
	Nonce   uint64
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
