package varint

import (
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

// Decode reads a variable bigInt from a stream.
func Decode(stream io.Reader) (*big.Int, error) {
	i := make([]byte, 1)

	_, err := stream.Read(i)
	if err != nil {
		return nil, err
	}

	if i[0] == 0xfd {
		n := make([]byte, 2)
		_, err := stream.Read(n)
		if err != nil {
			return nil, err
		}
		return endian.LittleEndianToBigInt(n), nil
	}

	if i[0] == 0xfe {
		n := make([]byte, 4)
		_, err := stream.Read(n)
		if err != nil {
			return nil, err
		}
		return endian.LittleEndianToBigInt(n), nil
	}

	if i[0] == 0xff {
		n := make([]byte, 8)
		_, err := stream.Read(n)
		if err != nil {
			return nil, err
		}
		return endian.LittleEndianToBigInt(n), nil
	}

	return endian.LittleEndianToBigInt(i), nil
}
