package bitcoin

import (
	"fmt"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

type Tx struct {
	Version int32
}

func NewTx(version int32) *Tx {
	return &Tx{version}
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

	_, err = ParseTxInputs(data)
	if err != nil {
		return nil, err
	}

	return NewTx(version), nil
}

func parseVersion(data io.Reader) (int32, error) {
	v := make([]byte, 4)

	_, err := data.Read(v)
	if err != nil {
		return 0, err
	}

	return endian.LittleEndianToInt32(v), nil
}
