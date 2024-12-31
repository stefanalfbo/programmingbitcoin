package bitcoin_test

import (
	"bytes"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestTx(t *testing.T) {
	t.Run("Parse version", func(t *testing.T) {
		stream := bytes.NewReader([]byte{0x01, 0x00, 0x00, 0x00, 0x01})
		tx, err := bitcoin.Parse(stream)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if tx.Version != 1 {
			t.Errorf("unexpected string: %s", tx.String())
		}
	})
}
