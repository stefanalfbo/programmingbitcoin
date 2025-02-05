package merkle

import (
	"encoding/binary"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type MerkleBlock struct {
	Version       int32
	PreviousBlock []byte
	MerkleRoot    []byte
	Timestamp     uint32
	Bits          []byte
	Nonce         uint32
	total         int
	hashes        [][]byte
	flags         []byte
}

func ParseMerkleBlock(reader io.Reader) (*MerkleBlock, error) {
	mb := MerkleBlock{}

	var version int32
	err := binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}
	mb.Version = version

	previousBlock := make([]byte, 32)
	err = binary.Read(reader, binary.LittleEndian, &previousBlock)
	if err != nil {
		return nil, err
	}
	mb.PreviousBlock = previousBlock

	merkleRoot := make([]byte, 32)
	err = binary.Read(reader, binary.LittleEndian, &merkleRoot)
	if err != nil {
		return nil, err
	}
	mb.MerkleRoot = merkleRoot

	var timestamp uint32
	err = binary.Read(reader, binary.LittleEndian, &timestamp)
	if err != nil {
		return nil, err
	}
	mb.Timestamp = timestamp

	bits := make([]byte, 4)
	err = binary.Read(reader, binary.LittleEndian, &bits)
	if err != nil {
		return nil, err
	}
	mb.Bits = bits

	var nonce uint32
	err = binary.Read(reader, binary.LittleEndian, &nonce)
	if err != nil {
		return nil, err
	}
	mb.Nonce = nonce

	var total uint32
	err = binary.Read(reader, binary.LittleEndian, &total)
	if err != nil {
		return nil, err
	}
	mb.total = int(total)

	hashCount, err := varint.Decode(reader)
	if err != nil {
		return nil, err
	}

	for i := uint64(0); i < hashCount; i++ {
		hash := make([]byte, 32)
		err = binary.Read(reader, binary.LittleEndian, &hash)
		if err != nil {
			return nil, err
		}
		mb.hashes = append(mb.hashes, hash)
	}

	flagsLength, err := varint.Decode(reader)
	if err != nil {
		return nil, err
	}

	flags := make([]byte, flagsLength)
	err = binary.Read(reader, binary.LittleEndian, &flags)
	if err != nil {
		return nil, err
	}
	mb.flags = flags

	return &mb, nil
}
