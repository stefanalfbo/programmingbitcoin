package merkle

import (
	"fmt"
	"math"
	"strings"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
)

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

type MerkleTree struct {
	total        int
	maxDepth     int
	Nodes        [][][]byte
	currentDepth int
	currentIndex int
}

func NewMerkleTree(total int) *MerkleTree {
	mt := MerkleTree{total: total}
	mt.maxDepth = int(math.Ceil(math.Log2(float64(total))))

	for depth := 0; depth < (mt.maxDepth + 1); depth++ {
		numItems := int32(math.Ceil(float64(mt.total) / math.Pow(2, float64(mt.maxDepth-depth))))

		levelHashes := make([][]byte, numItems)

		mt.Nodes = append(mt.Nodes, levelHashes)
	}

	mt.currentDepth = 0
	mt.currentIndex = 0

	return &mt
}

func (mt *MerkleTree) String() string {
	result := make([]string, len(mt.Nodes))

	for depth, level := range mt.Nodes {
		var short string
		var items []string
		for index, h := range level {
			if h == nil {
				short = "nil"
			} else {
				hex := fmt.Sprintf("%x", h)
				short = fmt.Sprintf("%s...", hex[:8])
			}

			if depth == mt.currentDepth && index == mt.currentIndex {
				items = append(items, fmt.Sprintf("*%s*", short[:8]))
			} else {
				items = append(items, short)
			}
		}
		result = append(result, fmt.Sprintf("depth %d: %s", depth, strings.Join(items, " ")))
	}
	return strings.Join(result, "\n")
}
