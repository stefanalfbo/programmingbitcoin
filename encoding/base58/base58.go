// Package base58 implements base58 encoding as used in Bitcoin.
package base58

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
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

func Decode(s string) ([]byte, error) {
	num := big.NewInt(0)
	for _, char := range s {
		num = num.Mul(num, big.NewInt(58))
		index := strings.Index(string(base58Alphabet), string(char))
		if index == -1 {
			return nil, fmt.Errorf("invalid character in BASE58")
		}
		num = num.Add(num, big.NewInt(int64(index)))
	}

	combined := make([]byte, 25)
	num.FillBytes(combined)
	checksum := combined[len(combined)-4:]
	withoutChecksum := combined[:len(combined)-4]

	calcChecksum := hash.Hash256(withoutChecksum)

	if !bytes.Equal(checksum, calcChecksum[:4]) {
		return nil, fmt.Errorf("checksum does not match")
	}

	return combined[1 : len(combined)-4], nil
}

func Checksum(data []byte) string {
	hash := hash.Hash256(data)
	dataWithChecksum := append(data, hash[:4]...)

	return Encode(dataWithChecksum)
}
