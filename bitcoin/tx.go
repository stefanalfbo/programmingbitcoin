package bitcoin

import (
	"fmt"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
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
	return []byte("dummy value")
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
