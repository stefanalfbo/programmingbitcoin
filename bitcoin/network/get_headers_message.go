package network

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type GetHeadersMessage struct {
	command    []byte
	version    uint32
	numHashes  uint8
	startBlock []byte
	endBlock   []byte
}

func NewGetHeadersMessage(version uint32, numHashes uint8, startBlock []byte, endBlock []byte) *GetHeadersMessage {
	command := []byte("getheaders")

	return &GetHeadersMessage{command, version, numHashes, startBlock, endBlock}
}

func (ghm *GetHeadersMessage) Command() []byte {
	return ghm.command
}

func (ghm *GetHeadersMessage) Serialize() ([]byte, error) {
	result := make([]byte, 0)

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, ghm.version)
	result = append(result, version...)

	numHashes, err := varint.Encode(big.NewInt(int64(ghm.numHashes)))
	if err != nil {
		return nil, err
	}
	result = append(result, numHashes...)

	startBlock := make([]byte, 32)
	slices.Reverse(ghm.startBlock)
	result = append(result, startBlock...)

	endBlock := make([]byte, 32)
	slices.Reverse(ghm.endBlock)
	result = append(result, endBlock...)

	return result, nil

}

func (ghm *GetHeadersMessage) Parse(data []byte) (Message, error) {
	return nil, fmt.Errorf("not implemented")
}
