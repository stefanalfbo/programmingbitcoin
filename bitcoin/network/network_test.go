package network_test

import (
	"bytes"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network"
)

func TestString(t *testing.T) {
	t.Run("to string, empty payload", func(t *testing.T) {
		command := []byte{0x76, 0x65, 0x72, 0x61, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		payload := []byte{}

		ne := network.NewNetworkEnvelope(command, payload, false)
		expected := "verack: "
		actual := ne.String()

		if expected != actual {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	})

	t.Run("to string, non-empty payload", func(t *testing.T) {
		command := []byte("hello, world")
		payload := []byte{0xDE, 0xAD, 0xBE, 0xEF}

		ne := network.NewNetworkEnvelope(command, payload, false)
		expected := "hello, world: deadbeef"
		actual := ne.String()

		if expected != actual {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	})
}

func TestParseNetworkEnvelope(t *testing.T) {
	t.Run("invalid network magic", func(t *testing.T) {
		message := bytes.NewReader([]byte{
			0xDE, 0xAD, 0xBE, 0xEF, // invalid magic
		})

		_, err := network.ParseNetworkEnvelope(message)
		if err == nil || err.Error() != "invalid magic: deadbeef" {
			t.Errorf("expected error but got: %v", err)
		}
	})

	t.Run("invalid checksum", func(t *testing.T) {
		message := bytes.NewReader([]byte{
			0xf9, 0xbe, 0xb4, 0xd9, // mainnet magic
			0x76, 0x65, 0x72, 0x61, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // command
			0x00, 0x00, 0x00, 0x00, // payload length
			0xDE, 0xAD, 0xBE, 0xEF, // invalid checksum
		})

		_, err := network.ParseNetworkEnvelope(message)
		if err == nil || err.Error() != "invalid checksum: deadbeef" {
			t.Errorf("expected error but got: %v", err)
		}
	})

	t.Run("empty payload", func(t *testing.T) {
		message := bytes.NewReader([]byte{
			0xf9, 0xbe, 0xb4, 0xd9, // mainnet magic
			0x76, 0x65, 0x72, 0x61, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // command
			0x00, 0x00, 0x00, 0x00, // payload length
			0x5d, 0xf6, 0xe0, 0xe2, // checksum
		})

		ne, err := network.ParseNetworkEnvelope(message)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := "verack: "
		if expected != ne.String() {
			t.Errorf("expected '%s' but got '%s'", expected, ne.String())
		}
	})

	t.Run("non-empty payload", func(t *testing.T) {
		message := bytes.NewReader([]byte{
			0xf9, 0xbe, 0xb4, 0xd9, // mainnet magic
			0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, // command
			0x65, 0x00, 0x00, 0x00, // payload length
			0x5f, 0x1a, 0x69, 0xd2, // checksum
			0x72, 0x11, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xbc, 0x8f, 0x5e, 0x54,
			0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff,
			0xc6, 0x1b, 0x64, 0x09, 0x20, 0x8d, 0x01, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xff, 0xff, 0xcb, 0x00, 0x71, 0xc0, 0x20, 0x8d,
			0x12, 0x80, 0x35, 0xcb, 0xc9, 0x79, 0x53, 0xf8,
			0x0f, 0x2f, 0x53, 0x61, 0x74, 0x6f, 0x73, 0x68,
			0x69, 0x3a, 0x30, 0x2e, 0x39, 0x2e, 0x33, 0x2f,
			0xcf, 0x05, 0x05, 0x00, 0x01, // payload
		})

		ne, err := network.ParseNetworkEnvelope(message)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := "version: 721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001"
		if expected != ne.String() {
			t.Errorf("expected '%s' but got '%s'", expected, ne.String())
		}
	})
}

func TestSerializeNetworkEnvelope(t *testing.T) {
	t.Run("empty payload", func(t *testing.T) {
		message := []byte{
			0xf9, 0xbe, 0xb4, 0xd9, // mainnet magic
			0x76, 0x65, 0x72, 0x61, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // command
			0x00, 0x00, 0x00, 0x00, // payload length
			0x5d, 0xf6, 0xe0, 0xe2, // checksum
		}

		ne, err := network.ParseNetworkEnvelope(bytes.NewReader(message))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		serialized := ne.Serialize()
		if !bytes.Equal(message, serialized) {
			t.Errorf("expected %x but got %x", message, serialized)
		}
	})

	t.Run("non-empty payload", func(t *testing.T) {
		message := []byte{
			0xf9, 0xbe, 0xb4, 0xd9, // mainnet magic
			0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, // command
			0x65, 0x00, 0x00, 0x00, // payload length
			0x5f, 0x1a, 0x69, 0xd2, // checksum
			0x72, 0x11, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xbc, 0x8f, 0x5e, 0x54,
			0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff,
			0xc6, 0x1b, 0x64, 0x09, 0x20, 0x8d, 0x01, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xff, 0xff, 0xcb, 0x00, 0x71, 0xc0, 0x20, 0x8d,
			0x12, 0x80, 0x35, 0xcb, 0xc9, 0x79, 0x53, 0xf8,
			0x0f, 0x2f, 0x53, 0x61, 0x74, 0x6f, 0x73, 0x68,
			0x69, 0x3a, 0x30, 0x2e, 0x39, 0x2e, 0x33, 0x2f,
			0xcf, 0x05, 0x05, 0x00, 0x01, // payload
		}

		ne, err := network.ParseNetworkEnvelope(bytes.NewReader(message))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		serialized := ne.Serialize()
		if !bytes.Equal(message, serialized) {
			t.Errorf("expected %x but got %x", message, serialized)
		}
	})
}
