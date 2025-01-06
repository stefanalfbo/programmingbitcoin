package bitcoin_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestParseScript(t *testing.T) {
	hexString := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	scriptPubKey, _ := hex.DecodeString(hexString)

	script, err := bitcoin.ParseScript(bytes.NewReader(scriptPubKey))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	instructionAsHex := script.Instructions()[0].Hex()
	expected := "304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601"
	if instructionAsHex != expected {
		t.Errorf("unexpected instruction: %s", instructionAsHex)
	}

	instructionAsHex = script.Instructions()[1].Hex()
	expected = "035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	if instructionAsHex != expected {
		t.Errorf("unexpected instruction: %s", instructionAsHex)
	}
}

func TestSerialize(t *testing.T) {
	hexString := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	scriptPubKey, _ := hex.DecodeString(hexString)

	script, _ := bitcoin.ParseScript(bytes.NewReader(scriptPubKey))

	serialized, err := script.Serialize()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	asHex := hex.EncodeToString(serialized)
	if asHex != hexString {
		t.Errorf("unexpected serialized script: %s", asHex)
	}
}

// def test_serialize(self):
// want = '6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937'
//         6b0247304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937
// script_pubkey = BytesIO(bytes.fromhex(want))
// script = Script.parse(script_pubkey)
// self.assertEqual(script.serialize().hex(), want)
