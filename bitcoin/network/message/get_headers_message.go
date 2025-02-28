package message

import (
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

// Return a headers packet containing the headers of blocks starting right
// after the last known hash in the block locator object, up to hash_stop or
// 2000 blocks, whichever comes first. To receive the next block headers, one
// needs to issue getheaders again with a new block locator object. Keep in
// mind that some clients may provide headers of blocks which are invalid if
// the block locator object contains a hash on the invalid branch.
type GetHeadersMessage struct {
	command []byte
	// The protocol version
	version uint32
	// Number of block locator hash entries
	hashCount uint64
	// Block locator object; newest back to genesis block (dense to start, but then sparse)
	startBlock [32]byte
	// Hash of the last desired block header; set to zero to get as many blocks as possible (2000)
	endBlock [32]byte
}

func NewGetHeadersMessage(version uint32, hashCount uint64, startBlock [32]byte, endBlock [32]byte) *GetHeadersMessage {
	command := []byte("getheaders")

	return &GetHeadersMessage{command, version, hashCount, startBlock, endBlock}
}

func (ghm *GetHeadersMessage) Command() []byte {
	return ghm.command
}

func (ghm *GetHeadersMessage) Serialize() ([]byte, error) {
	result := make([]byte, 0)

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, ghm.version)
	result = append(result, version...)

	hashCount, err := varint.Encode(uint64(ghm.hashCount))
	if err != nil {
		return nil, err
	}
	result = append(result, hashCount...)

	startBlock := [32]byte{}
	copy(startBlock[:], ghm.startBlock[:])
	slices.Reverse(startBlock[:])
	result = append(result, startBlock[:]...)

	endBlock := [32]byte{}
	copy(endBlock[:], ghm.endBlock[:])
	slices.Reverse(endBlock[:])
	result = append(result, endBlock[:]...)

	return result, nil

}

func (ghm *GetHeadersMessage) Parse(data []byte) (Message, error) {
	return nil, fmt.Errorf("not implemented")
}
