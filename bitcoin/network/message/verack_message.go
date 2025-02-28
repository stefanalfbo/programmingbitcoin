package message

import (
	"io"
)

// The verack message is sent in reply to version. This message consists of
// only a message header with the command string "verack". See:
// https://en.bitcoin.it/wiki/Protocol_documentation#verack
type VerAckMessage struct{}

// NewVerAckMessage returns a new VerAckMessage
func NewVerAckMessage() *VerAckMessage {
	return &VerAckMessage{}
}

// Command returns the command string which is used to determine which
// message is being sent.
func (vam *VerAckMessage) Command() []byte {
	return []byte("verack")
}

// Serialize serializes the message to a byte slice
func (vam *VerAckMessage) Serialize() ([]byte, error) {
	return []byte{}, nil
}

// Parse reads a VerAckMessage from a reader
func (vam *VerAckMessage) Parse(reader io.Reader) (Message, error) {
	return NewVerAckMessage(), nil
}
