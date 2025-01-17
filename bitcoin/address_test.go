package bitcoin_test

import (
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestH160PSHAddress(t *testing.T) {
	expectedMainnet := "2N3u1R6uwQfuobCqbCgBkpsgBxvr1tZpe7B"
	expectedTestnet := "2N3u1R6uwQfuobCqbCgBkpsgBxvr1tZpe7f"
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")

	// Mainnet
	got := bitcoin.H160ToP2SHAddress(h160, false)
	if got != expectedMainnet {
		t.Errorf("got %v, expected %v", got, expectedMainnet)
	}

	// Testnet
	got = bitcoin.H160ToP2SHAddress(h160, true)
	if got != expectedTestnet {
		t.Errorf("got %v, expected %v", got, expectedTestnet)
	}
}
