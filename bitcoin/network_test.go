package bitcoin_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestString(t *testing.T) {
	t.Run("to string, empty payload", func(t *testing.T) {
		command := []byte{0x76, 0x65, 0x72, 0x61, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		payload := []byte{}

		ne := bitcoin.NewNetworkEnvelope(command, payload, false)
		expected := "verack: "
		actual := ne.String()

		if expected != actual {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	})

	t.Run("to string, non-empty payload", func(t *testing.T) {
		command := []byte("hello, world!")
		payload := []byte{0xDE, 0xAD, 0xBE, 0xEF}

		ne := bitcoin.NewNetworkEnvelope(command, payload, false)
		expected := "hello, world!: deadbeef"
		actual := ne.String()

		if expected != actual {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	})
}
