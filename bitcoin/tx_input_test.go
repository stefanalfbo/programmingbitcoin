package bitcoin_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestParseInputs(t *testing.T) {
	setup := func() *bytes.Reader {
		hexString := "01813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

		dataBytes, _ := hex.DecodeString(hexString)

		return bytes.NewReader(dataBytes)
	}

	t.Run("String", func(t *testing.T) {
		stream := setup()

		inputs, err := bitcoin.ParseTxInputs(stream)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		txInput := inputs[0]

		if txInput.String() != "813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1:0" {
			t.Errorf("unexpected string: %s", txInput.String())
		}
	})

	t.Run("Parse inputs", func(t *testing.T) {
		stream := setup()

		inputs, err := bitcoin.ParseTxInputs(stream)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(inputs) != 1 {
			t.Errorf("unexpected number of inputs: %d", len(inputs))
		}

		txInput := inputs[0]

		if txInput.PrevIndex.Cmp(big.NewInt(0)) != 0 {
			t.Errorf("unexpected prev index: %d", txInput.PrevIndex)
		}

		if txInput.Sequence.Cmp(big.NewInt(0xfffffffe)) != 0 {
			t.Errorf("unexpected sequence: %d", txInput.Sequence)
		}
	})
}
