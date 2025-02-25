package bitcoin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
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

	return hash.Hash256(txSerialized)
}

func Parse(data io.Reader, isTestnet bool) (*Tx, error) {
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

	return NewTx(version, inputs, outputs, lockTime, isTestnet), nil
}

func parseVersion(data io.Reader) (int32, error) {
	version := make([]byte, 4)

	_, err := data.Read(version)
	if err != nil {
		return 0, err
	}

	return int32(binary.LittleEndian.Uint32(version)), nil
}

func parseLockTime(data io.Reader) (int32, error) {
	lockTime := make([]byte, 4)

	_, err := data.Read(lockTime)
	if err != nil {
		return 0, err
	}

	return int32(binary.LittleEndian.Uint32(lockTime)), nil
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
	result, err := varint.Encode(uint64(len(tx.Inputs)))
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
	result, err := varint.Encode(uint64(len(tx.Outputs)))
	if err != nil {
		return nil
	}

	for _, txOut := range tx.Outputs {
		result = append(result, txOut.Serialize()...)
	}

	return result
}

// Returns the fee of the transaction.
func (tx *Tx) Fee() (int64, error) {
	var inputSum int64 = 0
	var outputSum int64 = 0
	for _, txIn := range tx.Inputs {
		value, err := txIn.Value(tx.isTestnet)
		if err != nil {
			return 0, err
		}
		inputSum += int64(value)
	}
	for _, txOut := range tx.Outputs {
		outputSum += int64(txOut.Amount)
	}

	return inputSum - outputSum, nil
}

func (tx *Tx) SignatureHash(inputIndex int, redeemScript *Script) ([]byte, error) {
	signature := endian.BigIntToLittleEndian(big.NewInt(int64(tx.Version)), 4)

	length, err := varint.Encode(uint64(len(tx.Inputs)))
	if err != nil {
		return nil, err
	}

	signature = append(signature, length...)

	for i, txIn := range tx.Inputs {
		if i == inputIndex {
			if redeemScript != nil {
				signature, err = redeemScript.Serialize()
				if err != nil {
					return nil, err
				}
			} else {
				scriptSignature, err := txIn.ScriptPubKey(tx.isTestnet)
				if err != nil {
					return nil, err
				}
				tmpTxIn := NewTxInput(txIn.PrevTx, txIn.PrevIndex, scriptSignature, txIn.Sequence)

				signature = append(signature, tmpTxIn.Serialize()...)
			}
		} else {
			tmpTxIn := NewTxInput(txIn.PrevTx, txIn.PrevIndex, nil, txIn.Sequence)

			signature = append(signature, tmpTxIn.Serialize()...)
		}
	}

	outputLength, err := varint.Encode(uint64(len(tx.Outputs)))
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

	hashed := hash.Hash256(signature)

	return hashed, nil
}

func (tx *Tx) VerifyInput(inputIndex int) (bool, error) {
	txInput := tx.Inputs[inputIndex]
	scriptPubKey, err := txInput.ScriptPubKey(tx.isTestnet)
	if err != nil {
		return false, err
	}

	var redeemScript *Script
	if scriptPubKey.IsP2SHScriptPubKey() {
		instruction := scriptPubKey.instructions[len(scriptPubKey.instructions)-1]
		lenBytes, err := varint.Encode(uint64(instruction.Length()))
		if err != nil {
			return false, err
		}
		rawRedeemScript := append(lenBytes, instruction.Bytes()...)
		redeemScript, err = ParseScript(bytes.NewReader(rawRedeemScript))
		if err != nil {
			return false, err
		}
	} else {
		redeemScript = nil
	}

	z, err := tx.SignatureHash(inputIndex, redeemScript)
	if err != nil {
		return false, err
	}

	script := scriptPubKey.Add(txInput.ScriptSig)
	result, err := script.Evaluate(z)

	return result, err
}

// Verify this transaction
func (tx *Tx) Verify() (bool, error) {
	fee, err := tx.Fee()
	if err != nil {
		return false, err
	}

	if fee < 0 {
		return false, fmt.Errorf("Fee must be positive")
	}

	for index := 0; index < len(tx.Inputs); index++ {
		_, err := tx.VerifyInput(index)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (tx *Tx) SignInput(inputIndex int, privateKey *ecc.PrivateKey) (bool, error) {
	z, err := tx.SignatureHash(inputIndex, nil)
	if err != nil {
		return false, err
	}

	signature, err := privateKey.Sign(big.NewInt(0).SetBytes(z))
	if err != nil {
		return false, err
	}

	der := signature.DER()
	sig := append(der, byte(1))

	sec := privateKey.SECCompressed()

	sigInstruction, err := op.NewInstruction(sig)
	if err != nil {
		return false, err
	}

	secInstruction, err := op.NewInstruction(sec)
	if err != nil {
		return false, err
	}

	scriptSig := NewScript([]op.Instruction{*sigInstruction, *secInstruction})
	tx.Inputs[inputIndex].ScriptSig = scriptSig

	return tx.VerifyInput(inputIndex)
}

func (tx *Tx) IsCoinbase() bool {
	// Coinbase transactions have only one input
	if len(tx.Inputs) != 1 {
		return false
	}

	firstInput := tx.Inputs[0]

	// The previous transaction must be 32 bytes of 00
	zeroBytes := make([]byte, 32)
	if len(firstInput.PrevTx) != 32 || !bytes.Equal(firstInput.PrevTx, zeroBytes) {
		return false
	}

	// The previous index must be 0xffffffff
	if firstInput.PrevIndex.Cmp(big.NewInt(0xffffffff)) != 0 {
		return false
	}

	return true
}

func (tx *Tx) CoinbaseHeight() (int32, error) {
	if !tx.IsCoinbase() {
		return 0, fmt.Errorf("tx is not a coinbase transaction")
	}

	data := tx.Inputs[0].ScriptSig.instructions[0].Bytes()

	return int32(binary.LittleEndian.Uint32(data)), nil
}
