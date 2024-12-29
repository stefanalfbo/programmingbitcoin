package ecc

import (
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

func Hash256(msg string) *big.Int {
	h := sha256.New()
	_, err := h.Write([]byte(msg))
	if err != nil {
		return nil
	}
	firstRound := h.Sum(nil)
	h.Reset()

	_, err = h.Write(firstRound)
	if err != nil {
		return nil
	}

	return new(big.Int).SetBytes(h.Sum(nil))
}

// SHA256 followed by RIPEMD-160
func Hash160(data []byte) []byte {
	sha256Hasher := sha256.New()
	sha256Hasher.Write(data)
	sha256Hash := sha256Hasher.Sum(nil)

	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(sha256Hash)
	return ripemd160Hasher.Sum(nil)
}
