package merkle

import "github.com/stefanalfbo/programmingbitcoin/crypto/hash"

func Parent(left []byte, right []byte) []byte {
	concatenated := append(left, right...)

	return hash.Hash256(concatenated)
}
