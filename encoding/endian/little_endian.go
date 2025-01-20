package endian

import "math/big"

// LittleEndianToBigInt converts a little-endian byte slice to a big integer.
func LittleEndianToBigInt(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(reverseBytes(bytes))
}

func LittleEndianToInt32(bytes []byte) int32 {
	return int32(LittleEndianToBigInt(bytes).Int64())
}

func LittleEndianToUint32(bytes []byte) uint32 {
	return uint32(LittleEndianToBigInt(bytes).Int64())
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

func Int32ToLittleEndian(n int32) []byte {
	return BigIntToLittleEndian(big.NewInt(int64(n)), 4)
}

func Uint32ToLittleEndian(n uint32) []byte {
	return BigIntToLittleEndian(big.NewInt(int64(n)), 4)
}

func Uint64ToLittleEndian(n uint64) []byte {
	return BigIntToLittleEndian(big.NewInt(int64(n)), 8)
}

func Int64ToLittleEndian(n int64) []byte {
	return BigIntToLittleEndian(big.NewInt(n), 8)
}

func reverseBytes(bytes []byte) []byte {
	n := len(bytes)
	reversed := make([]byte, n)

	for i := 0; i < n; i++ {
		reversed[i] = bytes[n-1-i]
	}

	return reversed
}
