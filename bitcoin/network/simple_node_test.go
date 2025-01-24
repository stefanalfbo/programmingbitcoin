package network_test

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network"
)

func TestSimpleNode(t *testing.T) {
	port := 18333
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	defer listener.Close()

	nodeAddress := listener.Addr()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			_, _ = conn.Write([]byte("Echo: " + line))
		}
	}()

	t.Run("Send", func(t *testing.T) {
		simpleNode := network.NewSimpleNode(nodeAddress, false, false)
		versionMessage := network.NewVersionMessage()

		err := simpleNode.Send(versionMessage)
		if err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}
	})
}
