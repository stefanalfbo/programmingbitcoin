package bitcoin

import (
	"fmt"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type TxInput struct {
	PrevTx    []byte
	PrevIndex int
	ScriptSig []byte
	Sequence  int
}

func (txIn *TxInput) String() string {
	return fmt.Sprintf("TxInput: %x:%d", txIn.PrevTx, txIn.PrevIndex)
}

func ParseTxInputs(data io.Reader) ([]*TxInput, error) {
	_, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
