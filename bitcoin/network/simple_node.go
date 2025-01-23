package network

import (
	"fmt"
	"net"
)

type SimpleNode struct {
	host      string
	port      int
	isTestnet bool
	isLogging bool
}

func NewSimpleNode(host string, port int, isTestnet bool, isLogging bool) *SimpleNode {
	return &SimpleNode{host, port, isTestnet, isLogging}
}

// Send a message to the connected node
func (n *SimpleNode) Send(message Message) error {
	envelope := NewNetworkEnvelope(message.Command(), message.Serialize(), n.isTestnet)

	if n.isLogging {
		fmt.Println("Sending:", envelope)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", n.host, n.port))
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
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", n.host, n.port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data := make([]byte, 1024)
	_, err = conn.Read(data)
	if err != nil {
		return nil, err
	}

	envelope, err := ParseNetworkEnvelope(data)
	if err != nil {
		return nil, err
	}

	if n.isLogging {
		fmt.Println("Received:", envelope)
	}

	return envelope, nil
}
