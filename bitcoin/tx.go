package bitcoin

import (
	"fmt"
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type Tx struct {
	Version  int32
	Inputs   []*TxInput
	Outputs  []*TxOutput
	LockTime int32
}

func NewTx(version int32, inputs []*TxInput, outputs []*TxOutput, lockTime int32) *Tx {
	return &Tx{version, inputs, outputs, lockTime}
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

	return ecc.Hash256(string(txSerialized)).Bytes()
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

	return NewTx(version, inputs, outputs, lockTime), nil
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
