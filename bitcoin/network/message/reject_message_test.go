package message_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/network/message"
)

func TestParse(t *testing.T) {
	msg, _ := hex.DecodeString("07636f6d6d616e64110366656564726561736f6e0d6f7270686f6e20626c6f636b")

	rm := message.NewEmptyRejectMessage()
	rm.Parse(bytes.NewReader(msg))

	if string(rm.Message) != "command" {
		t.Errorf("expected message 'command', got %s", rm.Message)
	}

	if rm.CCode != message.REJECT_OBSOLETE {
		t.Errorf("expected ccode %d, got %d", message.REJECT_OBSOLETE, rm.CCode)
	}
}

func TestSerializeAndParseRejectMessage(t *testing.T) {
	var data [32]byte
	copy(data[:], []byte("orphan block"))
	rm := message.NewRejectMessage("command", message.REJECT_OBSOLETE, "reason", data)

	expected, _ := hex.DecodeString("07636f6d6d616e641106726561736f6e6f727068616e20626c6f636b0000000000000000000000000000000000000000")
	serialized, _ := rm.Serialize()

	if !bytes.Equal(serialized, expected) {
		t.Errorf("expected %x, got %x", expected, serialized)
	}

	rmParsed := message.NewEmptyRejectMessage()
	rmParsed.Parse(bytes.NewReader(serialized))

	if rm.Message != rmParsed.Message {
		t.Errorf("expected %s, got %s", rm.Message, rmParsed.Message)
	}

	if rm.CCode != rmParsed.CCode {
		t.Errorf("expected %d, got %d", rm.CCode, rmParsed.CCode)
	}

	if rm.Reason != rmParsed.Reason {
		t.Errorf("expected %s, got %s", rm.Reason, rmParsed.Reason)
	}

	if !bytes.Equal(rm.Data[:], rmParsed.Data[:]) {
		t.Errorf("expected %x, got %x", rm.Data, rmParsed.Data)
	}
}
