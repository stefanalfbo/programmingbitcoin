package network

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type HeadersMessage struct {
	command []byte
	blocks  []*bitcoin.Block
}

func NewHeadersMessage(blocks []*bitcoin.Block) *HeadersMessage {
	command := []byte("headers")

	return &HeadersMessage{command, blocks}
}

func (hm *HeadersMessage) Command() []byte {
	return hm.command
}

func (hm *HeadersMessage) Serialize() ([]byte, error) {
	result := make([]byte, 0)

	numBlocks, err := varint.Encode(big.NewInt(int64(len(hm.blocks))))
	if err != nil {
		return nil, err
	}
	result = append(result, numBlocks...)

	blocks := make([]byte, 0)
	for _, block := range hm.blocks {
		blockBytes, err := block.Serialize()
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, blockBytes...)
	}

	result = append(result, blocks...)

	return result, nil
}

func (hm *HeadersMessage) Parse(data []byte) (Message, error) {
	reader := bytes.NewReader(data)
	numHashes, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, err
	}
	blocks := make([]*bitcoin.Block, 0)
	for i := uint64(0); i < numHashes; i++ {
		blockData, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		block, err := bitcoin.ParseBlock(blockData)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return NewHeadersMessage(blocks), nil
}
