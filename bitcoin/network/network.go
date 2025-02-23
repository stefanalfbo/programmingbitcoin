package network

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
)

type NetworkMagic [4]byte

// Known network magic values
var (
	Mainnet  NetworkMagic = [4]byte{0xf9, 0xbe, 0xb4, 0xd9}
	Testnet3 NetworkMagic = [4]byte{0x0b, 0x11, 0x09, 0x07}
	Regtest  NetworkMagic = [4]byte{0xfa, 0xbf, 0xb5, 0xda}
	Signet   NetworkMagic = [4]byte{0x0a, 0x03, 0xcf, 0x40}
	Namecoin NetworkMagic = [4]byte{0xf9, 0xbe, 0xb4, 0xfe}
)

type Message interface {
	Command() []byte
	Serialize() ([]byte, error)
	Parse(io.Reader) (Message, error)
}

type NetworkEnvelope struct {
	// Magic value indicating message origin network, and used to seek to next
	// message when stream state is unknown
	magic []byte
	// ASCII string identifying the packet content, NULL padded (non-NULL padding
	// results in packet rejected)
	command []byte
	// The actual data
	payload []byte
}

func NewNetworkEnvelope(command []byte, payload []byte, isTestnet bool) *NetworkEnvelope {
	var magic []byte
	if isTestnet {
		magic = Testnet3[:]
	} else {
		magic = Mainnet[:]
	}

	trimmedCommand, err := trimCommand(bytes.NewReader(command))
	if err != nil {
		return nil
	}

	return &NetworkEnvelope{
		magic:   magic,
		command: trimmedCommand,
		payload: payload,
	}
}

func (ne *NetworkEnvelope) String() string {
	command := string(ne.command)
	payload := hex.EncodeToString(ne.payload)

	return fmt.Sprintf("%s: %s", command, payload)
}

func (ne *NetworkEnvelope) Command() []byte {
	return ne.command
}

func (ne *NetworkEnvelope) Payload() []byte {
	return ne.payload
}

func trimCommand(data io.Reader) ([]byte, error) {
	rawCommand := make([]byte, 12)
	_, err := data.Read(rawCommand)
	if err != nil {
		return nil, err
	}

	trimmed := make([]byte, 0)

	for _, b := range rawCommand {
		if b == 0x00 {
			break
		}

		trimmed = append(trimmed, b)
	}

	return trimmed, nil
}

func ParseNetworkEnvelope(data io.Reader) (*NetworkEnvelope, error) {
	magic := make([]byte, 4)
	_, err := data.Read(magic)
	if err != nil {
		return nil, err
	}

	if !isMainnet(magic) && !isTestnet(magic) {
		return nil, fmt.Errorf("invalid magic: %x", magic)
	}

	command, err := trimCommand(data)
	if err != nil {
		return nil, err
	}

	payloadLengthBytes := make([]byte, 4)
	_, err = data.Read(payloadLengthBytes)
	if err != nil {
		return nil, err
	}
	payloadLength := binary.LittleEndian.Uint32(payloadLengthBytes)

	checksum := make([]byte, 4)
	_, err = data.Read(checksum)
	if err != nil {
		return nil, err
	}

	var payload []byte
	if int(payloadLength) == 0 {
		payload = make([]byte, 0)
	} else {

		payload = make([]byte, payloadLength)
		_, err = data.Read(payload)
		if err != nil {
			return nil, err
		}
	}

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

	payloadLength := len(ne.payload)
	payloadLengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadLengthBytes, uint32(payloadLength))

	checksum := hash.Hash256(ne.payload)[:4]

	result := make([]byte, 0)
	result = append(result, ne.magic...)
	result = append(result, command...)
	result = append(result, payloadLengthBytes...)
	result = append(result, checksum...)
	result = append(result, ne.payload...)

	return result
}

func isMainnet(magic []byte) bool {
	return bytes.Equal(magic, Mainnet[:])
}

func isTestnet(magic []byte) bool {
	return bytes.Equal(magic, Testnet3[:])
}
