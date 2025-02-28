package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network/message"
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
func (n *SimpleNode) Send(msg message.Message) error {
	serializedMessage, err := msg.Serialize()
	if err != nil {
		return err
	}

	envelope := NewNetworkEnvelope(msg.Command(), serializedMessage, n.isTestnet)

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

func (n *SimpleNode) WaitFor(messages []message.Message) (*message.Message, error) {
	for {
		envelope, err := n.Read()
		if err != nil {
			return nil, err
		}

		for _, msg := range messages {

			if bytes.Equal(envelope.Command(), msg.Command()) {
				payload := bytes.NewReader(envelope.Payload())
				msg, err := msg.Parse(payload)
				if err != nil {
					return nil, err
				}

				return &msg, nil
			}
		}

		switch command := string(envelope.Command()); command {
		case "version":
			msg := message.NewVerAckMessage()
			n.Send(msg)
		case "ping":
			nonce := binary.LittleEndian.Uint64(envelope.Payload())
			message := message.NewPongMessage(nonce)
			n.Send(message)
		}
	}
}

func (n *SimpleNode) Handshake() error {
	msg := message.NewVersionMessage()
	err := n.Send(msg)
	if err != nil {
		return err
	}

	verAckMessage := message.NewVerAckMessage()
	_, err = n.WaitFor([]message.Message{verAckMessage})
	if err != nil {
		return err
	}

	return nil
}
