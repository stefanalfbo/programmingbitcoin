package bitcoin

import (
	"fmt"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

type Block struct {
	Version       int32
	PreviousBlock []byte
	MerkleRoot    []byte
	Timestamp     int32
	Bits          int32
	Nonce         int32
}

func NewBlock(version int32, previousBlock []byte, merkleRoot []byte, timestamp int32, bits int32, nonce int32) *Block {
	return &Block{version, previousBlock, merkleRoot, timestamp, bits, nonce}
}

func ParseBlock(data []byte) (*Block, error) {
	if len(data) < 80 {
		return nil, fmt.Errorf("data is too short")
	}

	version := endian.LittleEndianToInt32(data[:4])

	previousBlock := data[4:36]
	slices.Reverse(previousBlock)

	merkleRoot := data[36:68]
	slices.Reverse(merkleRoot)

	timestamp := endian.LittleEndianToInt32(data[68:72])
	bits := endian.LittleEndianToInt32(data[72:76])
	nonce := endian.LittleEndianToInt32(data[76:80])

	return &Block{
		Version:       version,
		PreviousBlock: previousBlock,
		MerkleRoot:    merkleRoot,
		Timestamp:     timestamp,
		Bits:          bits,
		Nonce:         nonce,
	}, nil
}

func (block *Block) Serialize() ([]byte, error) {
	data := make([]byte, 0)

	// Version
	data = append(data, endian.Int32ToLittleEndian(block.Version)...)
	// PreviousBlock
	previousBlock := make([]byte, 32)
	copy(previousBlock, block.PreviousBlock)
	slices.Reverse(previousBlock)
	data = append(data, block.PreviousBlock...)
	// MerkleRoot
	merkleRoot := make([]byte, 32)
	copy(merkleRoot, block.MerkleRoot)
	slices.Reverse(merkleRoot)
	data = append(data, block.MerkleRoot...)
	// Timestamp, Bits, Nonce
	data = append(data, endian.Int32ToLittleEndian(block.Timestamp)...)
	data = append(data, endian.Int32ToLittleEndian(block.Bits)...)
	data = append(data, endian.Int32ToLittleEndian(block.Nonce)...)

	return data, nil
}
