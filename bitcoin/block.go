package bitcoin

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

type Block struct {
	Version       int32
	PreviousBlock []byte
	MerkleRoot    []byte
	Timestamp     uint32
	Bits          []byte
	Nonce         uint32
}

func NewBlock(version int32, previousBlock []byte, merkleRoot []byte, timestamp uint32, bits []byte, nonce uint32) *Block {
	return &Block{version, previousBlock, merkleRoot, timestamp, bits, nonce}
}

func ParseBlock(data []byte) (*Block, error) {
	if len(data) < 80 {
		return nil, fmt.Errorf("data is too short")
	}

	version := endian.LittleEndianToInt32(data[:4])

	previousBlock := data[4:36]
	slices.Reverse(previousBlock)

	merkleRoot := data[36:68]
	slices.Reverse(merkleRoot)

	timestamp := endian.LittleEndianToUint32(data[68:72])
	bits := data[72:76] //endian.LittleEndianToUint32(data[72:76])
	nonce := endian.LittleEndianToUint32(data[76:80])

	return &Block{
		Version:       version,
		PreviousBlock: previousBlock,
		MerkleRoot:    merkleRoot,
		Timestamp:     timestamp,
		Bits:          bits,
		Nonce:         nonce,
	}, nil
}

func (block *Block) Serialize() ([]byte, error) {
	data := make([]byte, 0)

	// Version
	data = append(data, endian.Int32ToLittleEndian(block.Version)...)

	// PreviousBlock
	previousBlock := make([]byte, 32)
	copy(previousBlock, block.PreviousBlock)
	slices.Reverse(previousBlock)
	data = append(data, previousBlock...)

	// MerkleRoot
	merkleRoot := make([]byte, 32)
	copy(merkleRoot, block.MerkleRoot)
	slices.Reverse(merkleRoot)
	data = append(data, merkleRoot...)

	// Timestamp, Bits, Nonce
	data = append(data, endian.Uint32ToLittleEndian(block.Timestamp)...)
	data = append(data, block.Bits...)
	data = append(data, endian.Uint32ToLittleEndian(block.Nonce)...)

	return data, nil
}

// Returns the hash256 interpreted little endian of the block
func (block *Block) Hash() ([]byte, error) {
	serialized, err := block.Serialize()
	if err != nil {
		return nil, err
	}

	hashed := hash.Hash256(serialized)
	slices.Reverse(hashed)

	return hashed, nil
}

// Returns whether this block is signaling readiness for BIP9
func (block *Block) BIP9() bool {
	return block.Version>>29 == 0b001
}

func (block *Block) BIP91() bool {
	return block.Version>>4&1 == 1
}

func (block *Block) BIP141() bool {
	return block.Version>>1&1 == 1
}

func (block *Block) Target() *big.Int {
	return bitsToTarget(block.Bits)
}

func bitsToTarget(bits []byte) *big.Int {
	// target = coefficient * 256^(exponent-3)
	exponent := int64(bits[3])
	coefficient := big.NewInt(int64(endian.LittleEndianToUint32(bits[:3])))

	exponentPart := big.NewInt(0).Exp(big.NewInt(256), big.NewInt(exponent-3), nil)
	return big.NewInt(0).Mul(coefficient, exponentPart)
}

func (block *Block) Difficulty() *big.Int {
	// difficulty = 0xffff * 256^(0x1d - 3) / target
	lowest := big.NewInt(0).Mul(big.NewInt(0xffff), big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0x1d-3), nil))

	return big.NewInt(0).Div(lowest, block.Target())
}

func (block *Block) CheckProofOfWork() bool {
	serialized, err := block.Serialize()
	if err != nil {
		return false
	}

	sha := hash.Hash256(serialized)

	proof := endian.LittleEndianToBigInt(sha)
	target := block.Target()

	return proof.Cmp(target) == -1
}

func TargetToBits(target *big.Int) []byte {
	rawBytes := target.Bytes()

	var exponent int
	var coefficient []byte
	if rawBytes[0] > 0x7f {
		exponent = len(rawBytes) + 1
		coefficient = append([]byte{0x00}, rawBytes[:2]...)
	} else {
		exponent = len(rawBytes)
		coefficient = rawBytes[:3]
	}

	slices.Reverse(coefficient)

	bits := append(coefficient, byte(exponent))

	return bits
}

func CalculateNewBits(previousBits []byte, timeDifferential int64) []byte {
	TWO_WEEKS_IN_SECONDS := int64(60 * 60 * 24 * 14)

	// if the differential > 8 weeks, set to 8 weeks
	if timeDifferential > (TWO_WEEKS_IN_SECONDS * 4) {
		timeDifferential = TWO_WEEKS_IN_SECONDS * 4
	}

	// if the differential < 1/2 week, set to 1/2 week
	if timeDifferential < (TWO_WEEKS_IN_SECONDS / 4) {
		timeDifferential = TWO_WEEKS_IN_SECONDS / 4
	}

	target := bitsToTarget(previousBits)
	newTarget := new(big.Int).Mul(target, big.NewInt(timeDifferential))
	newTarget = new(big.Int).Div(newTarget, big.NewInt(TWO_WEEKS_IN_SECONDS))

	return TargetToBits(newTarget)
}
