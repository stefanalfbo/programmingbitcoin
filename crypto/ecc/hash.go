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
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil
	}

	sha256Hash := h.Sum(nil)

	ripemd160Hasher := ripemd160.New()
	_, err = ripemd160Hasher.Write(sha256Hash)
	if err != nil {
		return nil
	}

	return ripemd160Hasher.Sum(nil)
}
