package ecc_test

import (
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
)

func TestS256Field(t *testing.T) {
	t.Run("NewS256Field is valid", func(t *testing.T) {
		a, err := ecc.NewS256Field(big.NewInt(7))
		if err != nil {
			t.Errorf("NewS256Field: got error %v, expected nil", err)
		}
		if a == nil {
			t.Errorf("NewS256Field: got nil, expected valid S256Field")
		}
	})

	t.Run("NewS256Field is not valid", func(t *testing.T) {
		_, err := ecc.NewS256Field(big.NewInt(-14))

		if err == nil {
			t.Errorf("NewS256Field: expected error, got nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		a, _ := ecc.NewS256Field(big.NewInt(42))

		expected := "000000000000000000000000000000000000000000000000000000000000002a"
		if a.String() != expected {
			t.Errorf("String: got %v, expected %v", a.String(), expected)
		}
	})

	t.Run("not equal", func(t *testing.T) {
		a, _ := ecc.NewS256Field(big.NewInt(7))
		b, _ := ecc.NewS256Field(big.NewInt(6))

		if a.Equals(&b.FieldElement) {
			t.Errorf("Equals: got %v, expected %v", a, b)
		}
	})

	t.Run("equal", func(t *testing.T) {
		a, _ := ecc.NewS256Field(big.NewInt(7))

		if !a.Equals(&a.FieldElement) {
			t.Errorf("Equals: got %v, expected %v", a, a)
		}
	})
}
