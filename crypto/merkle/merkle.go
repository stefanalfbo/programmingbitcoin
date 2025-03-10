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

func (mt *MerkleTree) Up() {
	mt.currentDepth--
	mt.currentIndex /= 2
}

func (mt *MerkleTree) Left() {
	mt.currentDepth++
	mt.currentIndex *= 2
}

func (mt *MerkleTree) Right() {
	mt.currentDepth++
	mt.currentIndex = mt.currentIndex*2 + 1
}

func (mt *MerkleTree) Root() []byte {
	return mt.Nodes[0][0]
}

func (mt *MerkleTree) SetCurrentNode(value []byte) {
	mt.Nodes[mt.currentDepth][mt.currentIndex] = value
}

func (mt *MerkleTree) GetCurrentNode() []byte {
	return mt.Nodes[mt.currentDepth][mt.currentIndex]
}

func (mt *MerkleTree) GetLeftNode() []byte {
	return mt.Nodes[mt.currentDepth+1][mt.currentIndex*2]
}

func (mt *MerkleTree) GetRightNode() []byte {
	return mt.Nodes[mt.currentDepth+1][mt.currentIndex*2+1]
}

func (mt *MerkleTree) IsLeaf() bool {
	return mt.currentDepth == mt.maxDepth
}

func (mt *MerkleTree) RightExists() bool {
	return len(mt.Nodes[mt.currentDepth+1]) > mt.currentIndex*2+1
}

func (mt *MerkleTree) Populate(flagBits []bool, hashes [][]byte) {
	for {
		if mt.Root() != nil {
			break
		}

		if mt.IsLeaf() {
			flagBits = flagBits[:1]
			h := hashes[0]
			hashes = hashes[1:]
			mt.SetCurrentNode(h)
			mt.Up()
		} else {
			leftHash := mt.GetLeftNode()
			if leftHash == nil {
				flag := flagBits[0]
				flagBits = flagBits[1:]
				if flag {
					mt.Left()
				} else {
					h := hashes[0]
					hashes = hashes[1:]
					mt.SetCurrentNode(h)
					mt.Up()
				}
			} else if mt.RightExists() {
				rightHash := mt.GetRightNode()
				if rightHash == nil {
					mt.Right()
				} else {
					mt.SetCurrentNode(Parent(leftHash, rightHash))
					mt.Up()
				}
			} else {
				mt.SetCurrentNode(Parent(leftHash, leftHash))
				mt.Up()
			}
		}
	}

	if len(hashes) > 0 {
		panic("Not all hashes consumed")
	}

	for _, flag := range flagBits {
		if flag {
			panic("Not all flag bits are consumed")
		}
	}
}
