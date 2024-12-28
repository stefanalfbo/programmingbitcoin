// Package base58 implements base58 encoding as used in Bitcoin.
package base58

import (
	"crypto/sha256"
	"math/big"
)

var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Encode(s []byte) string {
	count := 0

	for _, c := range s {
		if c == 0 {
			count++
		} else {
			break
		}
	}

	prefix := make([]byte, count)
	for i := range prefix {
		prefix[i] = '1'
	}
	number := new(big.Int).SetBytes(s)
	encoded := make([]byte, 0)

	for number.Cmp(big.NewInt(0)) > 0 {
		n, r := new(big.Int).DivMod(number, big.NewInt(58), new(big.Int))
		encoded = append(encoded, base58Alphabet[r.Int64()])

		number = n
	}

	// Reverse the encoded bytes
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(append(prefix, encoded...))
}

func hash256(data []byte) []byte {
	h := sha256.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		return nil
	}
	firstRound := h.Sum(nil)
	h.Reset()

	_, err = h.Write(firstRound)
	if err != nil {
		return nil
	}

	return h.Sum(nil)
}

func Checksum(data []byte) string {
	hash := hash256(data)
	dataWithChecksum := append(data, hash[:4]...)

	return Encode(dataWithChecksum)
}
