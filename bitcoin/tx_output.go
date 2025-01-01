package bitcoin

import (
	"fmt"
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type TxOutput struct {
	Amount       *big.Int
	ScriptPubKey []byte
}

func (txOut *TxOutput) String() string {
	return fmt.Sprintf("%d:%d", txOut.Amount, txOut.ScriptPubKey)
}

func ParseTxOutputs(data io.Reader) ([]*TxOutput, error) {
	numberOfOutputs, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	outputs := make([]*TxOutput, numberOfOutputs.Int64())
	for i := 0; i < int(numberOfOutputs.Int64()); i++ {
		txOutput, err := parseTxOutput(data)
		if err != nil {
			return nil, err
		}

		outputs[i] = txOutput
	}

	return outputs, nil
}

func parseTxOutput(data io.Reader) (*TxOutput, error) {
	amount := make([]byte, 8)
	_, err := data.Read(amount)
	if err != nil {
		return nil, err
	}

	scriptPubKey, err := ParseScript(data)
	if err != nil {
		return nil, err
	}

	return &TxOutput{

		endian.LittleEndianToBigInt(amount),
		scriptPubKey,
	}, nil
}
