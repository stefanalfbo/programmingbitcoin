package bitcoin_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestParseOutputs(t *testing.T) {
	setup := func() *bytes.Reader {
		hexString := "02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

		dataBytes, _ := hex.DecodeString(hexString)

		return bytes.NewReader(dataBytes)
	}

	t.Run("String", func(t *testing.T) {
		stream := setup()

		outputs, err := bitcoin.ParseTxOutputs(stream)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		txOutput := outputs[0]

		if txOutput.String() != "32454049:notImplementedYet" {
			t.Errorf("unexpected string: %s", txOutput.String())
		}
	})

	t.Run("Parse outputs", func(t *testing.T) {
		stream := setup()

		outputs, err := bitcoin.ParseTxOutputs(stream)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(outputs) != 2 {
			t.Errorf("unexpected number of outputs: %d", len(outputs))
		}

		txOutput := outputs[0]

		if txOutput.Amount != 32_454_049 {
			t.Errorf("unexpected amount: %d", txOutput.Amount)
		}
	})
}
