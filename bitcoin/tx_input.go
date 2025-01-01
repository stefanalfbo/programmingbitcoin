package bitcoin

import (
	"fmt"
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type TxInput struct {
	PrevTx    []byte
	PrevIndex *big.Int
	ScriptSig []byte
	Sequence  *big.Int
}

func (txIn *TxInput) String() string {
	return fmt.Sprintf("%x:%d", txIn.PrevTx, txIn.PrevIndex)
}

func ParseTxInputs(data io.Reader) ([]*TxInput, error) {
	numberOfInputs, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	inputs := make([]*TxInput, numberOfInputs.Int64())
	for i := 0; i < int(numberOfInputs.Int64()); i++ {
		txInput, err := ParseTxInput(data)
		if err != nil {
			return nil, err
		}

		inputs[i] = txInput
	}

	return inputs, nil
}

func ParseTxInput(data io.Reader) (*TxInput, error) {
	previousTx := make([]byte, 32)
	_, err := data.Read(previousTx)
	if err != nil {
		return nil, err
	}

	previousTransactionIndex := make([]byte, 4)
	_, err = data.Read(previousTransactionIndex)
	if err != nil {
		return nil, err
	}

	scriptSignature, err := ParseScript(data)
	if err != nil {
		return nil, err
	}

	sequence := make([]byte, 4)
	_, err = data.Read(sequence)
	if err != nil {
		return nil, err
	}

	return &TxInput{
		previousTx,
		endian.LittleEndianToBigInt(previousTransactionIndex),
		scriptSignature,
		endian.LittleEndianToBigInt(sequence),
	}, nil
}
