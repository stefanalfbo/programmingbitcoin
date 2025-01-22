package network

type VerAckMessage struct {
	command []byte
}

func NewVerAckMessage() *VerAckMessage {
	command := []byte("verack")

	return &VerAckMessage{command: command}
}

func (vam *VerAckMessage) Command() []byte {
	return vam.command
}

func (vam *VerAckMessage) Serialize() []byte {
	return []byte{}
}
