package endian

import (
	"math/big"
	"slices"
)

// LittleEndianToBigInt converts a little-endian byte slice to a big integer.
func LittleEndianToBigInt(bytes []byte) *big.Int {
	slices.Reverse(bytes)
	return new(big.Int).SetBytes(bytes)
}

// BigIntToLittleEndian converts a big integer to a little-endian byte slice.
func BigIntToLittleEndian(n *big.Int, length int) []byte {
	bytes := n.Bytes()
	pad := length - len(bytes)

	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-i-1] = bytes[len(bytes)-i-1], bytes[i]
	}

	if pad > 0 {
		bytes = append(bytes, make([]byte, pad)...)
	}

	return bytes
}
