package bitcoin_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
)

func TestTx(t *testing.T) {
	setup := func() *bytes.Reader {
		hexString := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

		dataBytes, _ := hex.DecodeString(hexString)

		return bytes.NewReader(dataBytes)
	}

	t.Run("Parse version", func(t *testing.T) {
		stream := setup()
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if tx.Version != 1 {
			t.Errorf("unexpected string: %s", tx.String())
		}
	})

	t.Run("Parse LockTime", func(t *testing.T) {
		stream := setup()
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if tx.LockTime != 410393 {
			t.Errorf("unexpected LockTime, got: %v, expected: %v", tx.LockTime, 410393)
		}
	})

	t.Run("Fee", func(t *testing.T) {
		t.Skip("WIP")
		stream := setup()
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		fee, err := tx.Fee()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := big.NewInt(40_000)
		if fee.Cmp(expected) != 0 {
			t.Errorf("expected: %d, got: %d", expected, fee)
		}
	})

	t.Run("Parse inputs", func(t *testing.T) {
		t.Skip("WIP")
		stream := setup()
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(tx.Inputs) != 1 {
			t.Errorf("unexpected number of inputs: %d", len(tx.Inputs))
		}

		prevTxExpected, _ := hex.DecodeString("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
		if !bytes.Equal(tx.Inputs[0].PrevTx, prevTxExpected) {
			t.Errorf("expected: %x, got: %x", prevTxExpected, tx.Inputs[0].PrevTx)
		}

		expected := "813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1:0"
		if tx.Inputs[0].String() != expected {
			t.Errorf("expected: %s, got: %s", expected, tx.Inputs[0].String())
		}
	})

	t.Run("SignatureHash", func(t *testing.T) {
		t.Skip("WIP")
		expected := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006a47304402207db2402a3311a3b845b038885e3dd889c08126a8570f26a844e3e4049c482a11022010178cdca4129eacbeab7c44648bf5ac1f9cac217cd609d216ec2ebc8d242c0a012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
		hexString := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac194306003c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
		dataBytes, _ := hex.DecodeString(hexString)
		stream := bytes.NewReader(dataBytes)
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		z, err := tx.SignatureHash(0, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		privateKey, err := ecc.NewPrivateKey(big.NewInt(8675309))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		signature, err := privateKey.Sign(new(big.Int).SetBytes(z))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		der := signature.DER()
		sec := privateKey.SECCompressed()
		sig := append(der, 0x01)

		sigInstructions, _ := op.NewInstruction(sig)
		secInstructions, _ := op.NewInstruction(sec)
		instructions := []op.Instruction{*sigInstructions, *secInstructions}
		scriptSig := bitcoin.NewScript(instructions)

		tx.Inputs[0].ScriptSig = scriptSig

		value := hex.EncodeToString(tx.Serialize())
		if value != expected {
			t.Errorf("expected: %s, got: %s", expected, value)
		}
	})

	t.Run("IsCoinbase return false when tx is not a coinbase transaction", func(t *testing.T) {
		stream := setup()
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if tx.IsCoinbase() {
			t.Errorf("expected not a coinbase transaction")
		}
	})

	t.Run("IsCoinbase return true when tx is a coinbase transaction", func(t *testing.T) {
		hexString := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00ffffffff01faf20b58000000001976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac00000000"

		dataBytes, _ := hex.DecodeString(hexString)
		stream := bytes.NewReader(dataBytes)
		tx, err := bitcoin.Parse(stream, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !tx.IsCoinbase() {
			t.Errorf("expected a coinbase transaction")
		}
	})
}
