package network_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network"
)

func TestParse(t *testing.T) {
	message, _ := hex.DecodeString("07636f6d6d616e64110366656564726561736f6e0d6f7270686f6e20626c6f636b")

	rm := network.NewEmptyRejectMessage()
	rm.Parse(bytes.NewReader(message))

	if string(rm.Message) != "command" {
		t.Errorf("expected message 'command', got %s", rm.Message)
	}

	if rm.CCode != network.REJECT_OBSOLETE {
		t.Errorf("expected ccode %d, got %d", network.REJECT_OBSOLETE, rm.CCode)
	}

}
