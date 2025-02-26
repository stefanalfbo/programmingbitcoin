// This file contains the implementation of the hash functions
// used in the bitcoin codebase.
//
// The functions Hash256 and Hash160 are used to hash data with
// SHA-256 followed by SHA-256 and SHA-256 followed by RIPEMD-160 respectively.
//
// The function HashSHA1 is used to hash data with SHA-1 which is used in the
// bitcoin op code OP_SHA1.
//
// See https://en.bitcoin.it/wiki/Protocol_documentation#Hashes
package hash

import (
	"crypto/sha1"
	"crypto/sha256"

	//lint:ignore SA1019 we want to use the ripemd160 package, which is
	// used in the bitcoin codebase for creating bitcoin addresses
	"golang.org/x/crypto/ripemd160"
)

// SHA-256 followed by SHA-256
func Hash256(data []byte) []byte {
	h := sha256.New()
	_, err := h.Write(data)
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

func HashSHA1(data []byte) []byte {
	h := sha1.New()
	_, err := h.Write(data)
	if err != nil {
		return nil
	}
	return h.Sum(nil)
}
