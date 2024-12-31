package varint

import (
	"errors"
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

// Encodes a big integer to a varint
func Encode(n *big.Int) ([]byte, error) {
	if n.Cmp(big.NewInt(0xfd)) < 0 {
		return n.Bytes(), nil
	}

	if n.Cmp(big.NewInt(0x10000)) < 0 {
		return append([]byte{0xfd}, endian.BigIntToLittleEndian(n, 2)...), nil
	}

	if n.Cmp(big.NewInt(0x100000000)) < 0 {
		return append([]byte{0xfe}, endian.BigIntToLittleEndian(n, 4)...), nil
	}

	v, ok := new(big.Int).SetString("0x10000000000000000", 16)
	if ok && n.Cmp(v) < 0 {
		return append([]byte{0xff}, endian.BigIntToLittleEndian(n, 8)...), nil
	}

	return nil, errors.New("value too large")
}
