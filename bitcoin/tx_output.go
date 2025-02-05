package bitcoin

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type TxOutput struct {
	Amount       uint64
	ScriptPubKey Script
}

func (txOut *TxOutput) String() string {
	return fmt.Sprintf("%d:%s", txOut.Amount, txOut.ScriptPubKey.String())
}

func ParseTxOutputs(data io.Reader) ([]*TxOutput, error) {
	numberOfOutputs, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	outputs := make([]*TxOutput, numberOfOutputs)
	for i := 0; i < int(numberOfOutputs); i++ {
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
		binary.LittleEndian.Uint64(amount),
		*scriptPubKey,
	}, nil
}

// Returns the byte serialization of the transaction output.
func (txOut *TxOutput) Serialize() []byte {
	result, err := varint.Encode(txOut.Amount)
	if err != nil {
		return nil
	}

	// TODO: Implement the serialization of the ScriptPubKey.
	// result = append(result, txOut.ScriptPubKey.Serialize()...)

	return result
}
