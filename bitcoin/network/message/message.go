package message

import "io"

type Message interface {
	Command() []byte
	Serialize() ([]byte, error)
	Parse(io.Reader) (Message, error)
}
