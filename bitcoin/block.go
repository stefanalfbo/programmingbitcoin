package bitcoin

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/big"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
	"github.com/stefanalfbo/programmingbitcoin/crypto/merkle"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

type Block struct {
	Version       int32
	PreviousBlock [32]byte
	MerkleRoot    []byte
	Timestamp     uint32
	Bits          []byte
	Nonce         uint32
	txHashes      [][]byte
}

func NewBlock(version int32, previousBlock [32]byte, merkleRoot []byte, timestamp uint32, bits []byte, nonce uint32, txHashes [][]byte) *Block {
	return &Block{version, previousBlock, merkleRoot, timestamp, bits, nonce, txHashes}
}

func ParseBlock(data io.Reader) (*Block, error) {
	versionBytes := make([]byte, 4)
	_, err := data.Read(versionBytes)
	if err != nil {
		return nil, err
	}
	version := int32(binary.LittleEndian.Uint32(versionBytes))

	previousBlockSlice := make([]byte, 32)
	_, err = data.Read(previousBlockSlice)
	if err != nil {
		return nil, err
	}
	slices.Reverse(previousBlockSlice)
	var previousBlock [32]byte
	copy(previousBlock[:], previousBlockSlice)

	merkleRoot := make([]byte, 32)
	_, err = data.Read(merkleRoot)
	if err != nil {
		return nil, err
	}
	slices.Reverse(merkleRoot)

	timeBytes := make([]byte, 4)
	_, err = data.Read(timeBytes)
	if err != nil {
		return nil, err
	}
	timestamp := binary.LittleEndian.Uint32(timeBytes)

	bits := make([]byte, 4)
	_, err = data.Read(bits)
	if err != nil {
		return nil, err
	}

	nonceBytes := make([]byte, 4)
	_, err = data.Read(nonceBytes)
	if err != nil {
		return nil, err
	}
	nonce := binary.LittleEndian.Uint32(nonceBytes)

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
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(block.Version))
	data = append(data, version...)

	// PreviousBlock
	previousBlock := make([]byte, 32)
	copy(previousBlock, block.PreviousBlock[:])
	slices.Reverse(previousBlock)
	data = append(data, previousBlock...)

	// MerkleRoot
	merkleRoot := make([]byte, 32)
	copy(merkleRoot, block.MerkleRoot)
	slices.Reverse(merkleRoot)
	data = append(data, merkleRoot...)

	// Timestamp, Bits, Nonce
	timestamp := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestamp, block.Timestamp)
	data = append(data, timestamp...)

	data = append(data, block.Bits...)

	nonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonce, block.Nonce)
	data = append(data, nonce...)

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

func (block *Block) ValidateMerkleRoot() bool {
	hashes := make([][]byte, 0)

	for _, txHash := range block.txHashes {
		hash := make([]byte, 32)
		copy(hash, txHash)
		slices.Reverse(hash)
		hashes = append(hashes, hash)
	}

	root := merkle.Root(hashes)

	slices.Reverse(root)

	return bytes.Equal(root, block.MerkleRoot)
}
