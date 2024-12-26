package ecc_test

import (
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestSignature(t *testing.T) {
	r := big.NewInt(12345)
	s := big.NewInt(67890)

	t.Run("NewSignature", func(t *testing.T) {
		sig := ecc.NewSignature(r, s)

		if sig == nil {
			t.Errorf("NewSignature: got nil, expected non-nil")
			return
		}
		if sig.R.Cmp(r) != 0 {
			t.Errorf("NewSignature: got r %v, expected %v", sig.R, r)
		}
		if sig.S.Cmp(s) != 0 {
			t.Errorf("NewSignature: got s %v, expected %v", sig.S, s)
		}
	})

	t.Run("String", func(t *testing.T) {
		sig := ecc.NewSignature(r, s)
		expected := "Signature(3039, 10932)"

		if sig.String() != expected {
			t.Errorf("String: got %v, expected %v", sig.String(), expected)
		}
	})
}
