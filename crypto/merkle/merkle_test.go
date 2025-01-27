package merkle_test

import (
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/merkle"
)

func TestParent(t *testing.T) {
	expectedParent := "8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd"
	left, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	right, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")

	parent := merkle.Parent(left, right)

	parentHex := hex.EncodeToString(parent)
	if parentHex != expectedParent {
		t.Errorf("Parent was incorrect, got: %s, want: %s.", parentHex, expectedParent)
	}
}
