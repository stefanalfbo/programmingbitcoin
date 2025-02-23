package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type SimpleNode struct {
	address   net.Addr
	isTestnet bool
	isLogging bool
}

func NewSimpleNode(address net.Addr, isTestnet bool, isLogging bool) *SimpleNode {
	return &SimpleNode{address, isTestnet, isLogging}
}

// Send a message to the connected node
func (n *SimpleNode) Send(message Message) error {
	serializedMessage, err := message.Serialize()
	if err != nil {
		return err
	}

	envelope := NewNetworkEnvelope(message.Command(), serializedMessage, n.isTestnet)

	if n.isLogging {
		fmt.Println("Sending:", envelope)
	}

	conn, err := net.Dial("tcp", n.address.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(envelope.Serialize())
	if err != nil {
		return err
	}

	return nil
}

// Read a message from the socket
func (n *SimpleNode) Read() (*NetworkEnvelope, error) {
	conn, err := net.Dial("tcp", n.address.String())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data := make([]byte, 1024)
	_, err = conn.Read(data)
	if err != nil {
		return nil, err
	}

	envelope, err := ParseNetworkEnvelope(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if n.isLogging {
		fmt.Println("Received:", envelope)
	}

	return envelope, nil
}

func (n *SimpleNode) WaitFor(messages []Message) (*Message, error) {
	for {
		envelope, err := n.Read()
		if err != nil {
			return nil, err
		}

		for _, message := range messages {

			if bytes.Equal(envelope.Command(), message.Command()) {
				payload := bytes.NewReader(envelope.Payload())
				msg, err := message.Parse(payload)
				if err != nil {
					return nil, err
				}

				return &msg, nil
			}
		}

		switch command := string(envelope.Command()); command {
		case "version":
			message := NewVerAckMessage()
			n.Send(message)
		case "ping":
			nonce := binary.LittleEndian.Uint64(envelope.Payload())
			message := NewPongMessage(nonce)
			n.Send(message)
		}
	}
}

func (n *SimpleNode) Handshake() error {
	message := NewVersionMessage()
	err := n.Send(message)
	if err != nil {
		return err
	}

	_, err = n.WaitFor([]Message{NewVerAckMessage()})
	if err != nil {
		return err
	}

	return nil
}
