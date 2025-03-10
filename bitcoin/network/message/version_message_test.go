package message_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network/message"
)

func TestSerialize(t *testing.T) {
	t.Run("serialize version message", func(t *testing.T) {
		expectedAsHexString := "7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000"

		vm := message.NewVersionMessage()
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

	t.Run("parse version message", func(t *testing.T) {
		msg, _ := hex.DecodeString("7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000")

		vm := message.NewVersionMessage()
		vm.Parse(bytes.NewReader(msg))

		if vm.Version != 70015 {
			t.Errorf("expected version 70015, got %d", vm.Version)
		}

		if vm.Services != 0 {
			t.Errorf("expected services 0, got %d", vm.Services)
		}

		if vm.Timestamp != 0 {
			t.Errorf("expected timestamp 0, got %d", vm.Timestamp)
		}

		actual, err := vm.Serialize()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !bytes.Equal(actual, msg) {
			t.Errorf("expected '%x' but got '%x'", msg, actual)
		}
	})
}
