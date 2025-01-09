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
	ScriptSig *Script
	Sequence  *big.Int
}

func NewTxInput(prevTx []byte, prevIndex *big.Int, scriptSig *Script, sequence *big.Int) *TxInput {
	return &TxInput{
		prevTx,
		prevIndex,
		scriptSig,
		sequence,
	}
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
		txInput, err := parseTxInput(data)
		if err != nil {
			return nil, err
		}

		inputs[i] = txInput
	}

	return inputs, nil
}

func parseTxInput(data io.Reader) (*TxInput, error) {
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

// Returns the byte serialization of the transaction input.
func (txIn *TxInput) Serialize() []byte {
	result := make([]byte, 0)

	result = append(result, txIn.PrevTx...)
	result = append(result, endian.BigIntToLittleEndian(txIn.PrevIndex, 4)...)
	// TODO: Implement the serialization of the ScriptSig.
	// result = append(result, txIn.ScriptSig.Serialize()...)
	result = append(result, endian.BigIntToLittleEndian(txIn.Sequence, 4)...)

	return result
}

func (txIn *TxInput) fetchTransaction(isTestnet bool) (*Tx, error) {
	txFetcher := NewTxFetcher(MemPoolFetcher, isTestnet)

	previousTxHex := fmt.Sprintf("%x", txIn.PrevTx)
	tx, err := txFetcher.Fetch(previousTxHex, false)
	if err != nil {
		return nil, err
	}

	return tx, nil

}

// Get the output value by looking up tx hash. Returns the amount in satoshi.
func (txIn *TxInput) Value(isTestnet bool) (*big.Int, error) {
	tx, err := txIn.fetchTransaction(isTestnet)
	if err != nil {
		return nil, err
	}

	return tx.Outputs[txIn.PrevIndex.Int64()].Amount, nil
}

// Get the ScriptPubKey by looking up the tx hash. Returns a Script object.
func (txIn *TxInput) ScriptPubKey(isTestnet bool) (*Script, error) {
	tx, err := txIn.fetchTransaction(isTestnet)
	if err != nil {
		return nil, err
	}

	return &tx.Outputs[txIn.PrevIndex.Int64()].ScriptPubKey, nil
}
