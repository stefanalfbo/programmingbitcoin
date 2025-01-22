package network

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

var mainnetMagic = []byte{0xf9, 0xbe, 0xb4, 0xd9}
var testnetMagic = []byte{0x0b, 0x11, 0x09, 0x07}

type Message interface {
	Command() []byte
	Serialize() []byte
}

type NetworkEnvelope struct {
	magic   []byte
	command []byte
	payload []byte
}

func NewNetworkEnvelope(command []byte, payload []byte, isTestnet bool) *NetworkEnvelope {
	var magic []byte
	if isTestnet {
		magic = testnetMagic
	} else {
		magic = mainnetMagic
	}

	return &NetworkEnvelope{
		magic:   magic,
		command: command,
		payload: payload,
	}
}

func (ne *NetworkEnvelope) String() string {
	command := string(trimCommand(ne.command))
	payload := hex.EncodeToString(ne.payload)

	return fmt.Sprintf("%s: %s", command, payload)
}

func trimCommand(command []byte) []byte {
	trimmed := make([]byte, 0)

	for _, b := range command {
		if b == 0x00 {
			break
		}

		trimmed = append(trimmed, b)
	}

	return trimmed
}

func ParseNetworkEnvelope(data []byte) (*NetworkEnvelope, error) {
	magic := data[:4]

	if !bytes.Equal(magic, mainnetMagic) && !bytes.Equal(magic, testnetMagic) {
		return nil, fmt.Errorf("invalid magic: %x", magic)
	}

	command := trimCommand(data[4:16])
	payloadLength := endian.LittleEndianToInt32(data[16:20])
	checksum := data[20:24]
	payload := data[24:]

	if len(payload) != int(payloadLength) {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	calculatedChecksum := hash.Hash256(payload)
	if !bytes.Equal(checksum, calculatedChecksum[:4]) {
		return nil, fmt.Errorf("invalid checksum: %x", checksum)
	}

	return &NetworkEnvelope{
		magic:   magic,
		command: command,
		payload: payload,
	}, nil
}

func (ne *NetworkEnvelope) Serialize() []byte {
	command := make([]byte, 12)
	copy(command, ne.command)

	payloadLength := endian.Int32ToLittleEndian(int32(len(ne.payload)))
	checksum := hash.Hash256(ne.payload)[:4]

	result := make([]byte, 0)
	result = append(result, ne.magic...)
	result = append(result, command...)
	result = append(result, payloadLength...)
	result = append(result, checksum...)
	result = append(result, ne.payload...)

	return result
}
