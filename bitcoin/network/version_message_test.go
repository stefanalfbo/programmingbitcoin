package network_test

import (
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network"
)

func TestSerialize(t *testing.T) {
	t.Run("serialize version message", func(t *testing.T) {
		expectedAsHexString := "7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000"

		vm := network.NewVersionMessage()
		vm.Nonce = 0

		actual, err := vm.Serialize()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		actualAsHexString := hex.EncodeToString(actual)
		if expectedAsHexString != actualAsHexString {
			t.Errorf("expected '%s' but got '%s'", expectedAsHexString, actualAsHexString)
		}
	})
}
