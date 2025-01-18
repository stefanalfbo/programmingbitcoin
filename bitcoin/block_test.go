package bitcoin_test

import (
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestBlock(t *testing.T) {

	t.Run("ParseBlock return block with correct values", func(t *testing.T) {
		hexString := "000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if block.Version != 536870912 {
			t.Errorf("expected version 536870912, got %d", block.Version)
		}

		if block.PreviousBlock == nil {
			t.Errorf("expected previous block not to be nil")
		}

		if block.MerkleRoot == nil {
			t.Errorf("expected merkle root not to be nil")
		}

		if block.Timestamp != 1496586576 {
			t.Errorf("expected timestamp 1496586576, got %d", block.Timestamp)
		}

		if block.Bits != 402774100 {
			t.Errorf("expected bits 402774100, got %d", block.Bits)
		}

		if block.Nonce != 773550753 {
			t.Errorf("expected nonce 773550753, got %d", block.Nonce)
		}
	})
}
