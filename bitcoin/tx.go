package bitcoin

import (
	"fmt"
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type Tx struct {
	Version   int32
	Inputs    []*TxInput
	Outputs   []*TxOutput
	LockTime  int32
	isTestnet bool
}

func NewTx(version int32, inputs []*TxInput, outputs []*TxOutput, lockTime int32, isTestnet bool) *Tx {
	return &Tx{version, inputs, outputs, lockTime, isTestnet}
}

func (tx *Tx) String() string {
	return fmt.Sprintf("Tx: %s, version: %d", tx.Id(), tx.Version)
}

// Human-readable hexadecimal of the transaction hash.
func (tx *Tx) Id() string {
	return fmt.Sprintf("%x", tx.hash())
}

// Binary hash of the legacy serialization.
func (tx *Tx) hash() []byte {
	txSerialized := tx.Serialize()

	return hash.Hash256(txSerialized).Bytes()
}

func Parse(data io.Reader) (*Tx, error) {
	version, err := parseVersion(data)
	if err != nil {
		return nil, err
	}

	inputs, err := ParseTxInputs(data)
	if err != nil {
		return nil, err
	}

	outputs, err := ParseTxOutputs(data)
	if err != nil {
		return nil, err
	}

	lockTime, err := parseLockTime(data)
	if err != nil {
		return nil, err
	}

	return NewTx(version, inputs, outputs, lockTime, false), nil
}

func parseVersion(data io.Reader) (int32, error) {
	version := make([]byte, 4)

	_, err := data.Read(version)
	if err != nil {
		return 0, err
	}

	return endian.LittleEndianToInt32(version), nil
}

func parseLockTime(data io.Reader) (int32, error) {
	lockTime := make([]byte, 4)

	_, err := data.Read(lockTime)
	if err != nil {
		return 0, err
	}

	return endian.LittleEndianToInt32(lockTime), nil
}

// Returns the byte serialization of the transaction.
func (tx *Tx) Serialize() []byte {
	result := endian.BigIntToLittleEndian(big.NewInt(int64(tx.Version)), 4)
	result = append(result, tx.serializeInputs()...)
	result = append(result, tx.serializeOutputs()...)
	result = append(result, endian.BigIntToLittleEndian(big.NewInt(int64(tx.LockTime)), 4)...)

	return result
}

// Returns the byte serialization of the transaction inputs.
func (tx *Tx) serializeInputs() []byte {
	result, err := varint.Encode(big.NewInt(int64(len(tx.Inputs))))
	if err != nil {
		return nil
	}

	for _, txIn := range tx.Inputs {
		result = append(result, txIn.Serialize()...)
	}

	return result
}

// Returns the byte serialization of the transaction outputs.
func (tx *Tx) serializeOutputs() []byte {
	result, err := varint.Encode(big.NewInt(int64(len(tx.Outputs))))
	if err != nil {
		return nil
	}

	for _, txOut := range tx.Outputs {
		result = append(result, txOut.Serialize()...)
	}

	return result
}

// Returns the fee of the transaction.
func (tx *Tx) Fee(testnet bool) (*big.Int, error) {
	inputSum, outputSum := big.NewInt(0), big.NewInt(0)
	for _, txIn := range tx.Inputs {
		value, err := txIn.Value(testnet)
		if err != nil {
			return nil, err
		}
		inputSum = new(big.Int).Add(value, inputSum)
	}
	for _, txOut := range tx.Outputs {
		outputSum = new(big.Int).Add(txOut.Amount, outputSum)
	}

	return new(big.Int).Sub(inputSum, outputSum), nil
}

func (tx *Tx) SignatureHash(inputIndex int) ([]byte, error) {
	signature := endian.BigIntToLittleEndian(big.NewInt(int64(tx.Version)), 4)

	length, err := varint.Encode(big.NewInt(int64(len(tx.Inputs))))
	if err != nil {
		return nil, err
	}

	signature = append(signature, length...)

	for i, txIn := range tx.Inputs {
		if i == inputIndex {
			scriptSignature, err := txIn.ScriptPubKey(tx.isTestnet)
			if err != nil {
				return nil, err
			}
			tmpTxIn := NewTxInput(txIn.PrevTx, txIn.PrevIndex, scriptSignature, txIn.Sequence)

			signature = append(signature, tmpTxIn.Serialize()...)
		} else {
			tmpTxIn := NewTxInput(txIn.PrevTx, txIn.PrevIndex, nil, txIn.Sequence)

			signature = append(signature, tmpTxIn.Serialize()...)
		}
	}

	outputLength, err := varint.Encode(big.NewInt(int64(len(tx.Outputs))))
	if err != nil {
		return nil, err
	}

	signature = append(signature, outputLength...)
	for _, txOut := range tx.Outputs {
		signature = append(signature, txOut.Serialize()...)
	}

	lockTime := endian.BigIntToLittleEndian(big.NewInt(int64(tx.LockTime)), 4)
	signature = append(signature, lockTime...)

	SIGHASH_ALL := 1
	signature = append(signature, endian.BigIntToLittleEndian(big.NewInt(int64(SIGHASH_ALL)), 4)...)

	hashed := hash.Hash256(signature).Bytes()

	return hashed, nil
}
