package merkle

import "github.com/stefanalfbo/programmingbitcoin/crypto/hash"

// Takes the binary hashes and calculates the hash256
func Parent(left []byte, right []byte) []byte {
	concatenated := append(left, right...)

	return hash.Hash256(concatenated)
}

// Takes a list of binary hashes and returns a list that's half the length
func ParentLevel(level [][]byte) [][]byte {
	// if the number of hashes is odd, duplicate the last item
	if len(level)%2 != 0 {
		level = append(level, level[len(level)-1])
	}

	parentLevel := make([][]byte, len(level)/2)

	for i := 0; i < len(level); i += 2 {
		parentLevel[i/2] = Parent(level[i], level[i+1])
	}

	return parentLevel
}

// Takes a list of binary hashes and returns the merkle root
func Root(data [][]byte) []byte {
	level := data

	for len(level) > 1 {
		level = ParentLevel(level)
	}

	return level[0]
}
