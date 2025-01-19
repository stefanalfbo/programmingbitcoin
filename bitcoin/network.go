package bitcoin

import (
	"encoding/hex"
	"fmt"
)

type NetworkEnvelope struct {
	magic   []byte
	command []byte
	payload []byte
}

func NewNetworkEnvelope(command []byte, payload []byte, isTestnet bool) *NetworkEnvelope {
	var magic []byte
	if isTestnet {
		magic = []byte{0x0b, 0x11, 0x09, 0x07}
	} else {
		magic = []byte{0xf9, 0xbe, 0xb4, 0xd9}
	}

	return &NetworkEnvelope{
		magic:   magic,
		command: command,
		payload: payload,
	}
}

func (ne *NetworkEnvelope) String() string {
	trimmed := make([]byte, 0)

	for _, b := range ne.command {
		if b == 0x00 {
			break
		}

		trimmed = append(trimmed, b)
	}

	command := string(trimmed)
	payload := hex.EncodeToString(ne.payload)

	return fmt.Sprintf("%s: %s", command, payload)
}
