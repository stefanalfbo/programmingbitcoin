package ecc_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
)

func TestSignature(t *testing.T) {
	t.Run("NewSignature", func(t *testing.T) {
		r := big.NewInt(12345)
		s := big.NewInt(67890)
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
		r := big.NewInt(12345)
		s := big.NewInt(67890)
		sig := ecc.NewSignature(r, s)
		expected := "Signature(3039, 10932)"

		if sig.String() != expected {
			t.Errorf("String: got %v, expected %v", sig.String(), expected)
		}
	})

	t.Run("DER", func(t *testing.T) {
		expected := "3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec"
		r := new(big.Int)
		r.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
		s := new(big.Int)
		s.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
		signature := ecc.NewSignature(r, s)

		der := signature.DER()
		derHex := fmt.Sprintf("%x", der)

		if derHex != expected {
			t.Errorf("DER: got %s, expected %s", derHex, expected)
		}
	})

	t.Run("ParseSignature", func(t *testing.T) {
		testCases := []struct {
			r *big.Int
			s *big.Int
		}{
			{big.NewInt(1), big.NewInt(2)},
			{big.NewInt(42), big.NewInt(1337)},
			{big.NewInt(7777), big.NewInt(8888)},
		}

		for _, tc := range testCases {
			signature := ecc.NewSignature(tc.r, tc.s)
			der := signature.DER()

			signature2, err := ecc.ParseDER(der)
			if err != nil {
				t.Fatalf("ParseSignature: unexpected error: %v", err)
			}

			if signature2.R.Cmp(tc.r) != 0 {
				t.Errorf("ParseSignature: got r %v, expected %v", signature2.R, tc.r)
			}

			if signature2.S.Cmp(tc.s) != 0 {
				t.Errorf("ParseSignature: got s %v, expected %v", signature2.S, tc.s)
			}
		}
	})
}
