package bitcoin

import (
	"fmt"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

type Block struct {
	Version       int32
	PreviousBlock []byte
	MerkleRoot    []byte
	Timestamp     uint32
	Bits          uint32
	Nonce         uint32
}

func NewBlock(version int32, previousBlock []byte, merkleRoot []byte, timestamp uint32, bits uint32, nonce uint32) *Block {
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

	timestamp := endian.LittleEndianToUint32(data[68:72])
	bits := endian.LittleEndianToUint32(data[72:76])
	nonce := endian.LittleEndianToUint32(data[76:80])

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
	data = append(data, previousBlock...)

	// MerkleRoot
	merkleRoot := make([]byte, 32)
	copy(merkleRoot, block.MerkleRoot)
	slices.Reverse(merkleRoot)
	data = append(data, merkleRoot...)

	// Timestamp, Bits, Nonce
	data = append(data, endian.Uint32ToLittleEndian(block.Timestamp)...)
	data = append(data, endian.Uint32ToLittleEndian(block.Bits)...)
	data = append(data, endian.Uint32ToLittleEndian(block.Nonce)...)

	return data, nil
}

// Returns the hash256 interpreted little endian of the block
func (block *Block) Hash() ([]byte, error) {
	serialized, err := block.Serialize()
	if err != nil {
		return nil, err
	}

	hashed := hash.Hash256(serialized)
	slices.Reverse(hashed)

	return hashed, nil
}

// Returns whether this block is signaling readiness for BIP9
func (block *Block) BIP9() bool {
	return block.Version>>29 == 0b001
}

func (block *Block) BIP91() bool {
	return block.Version>>4&1 == 1
}
