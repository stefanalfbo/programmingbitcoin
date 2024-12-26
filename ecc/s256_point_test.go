package ecc_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestS256Point(t *testing.T) {
	t.Run("New S256 point", func(t *testing.T) {
		x, _ := ecc.NewS256Field(ecc.Secp256k1.Gx)
		y, _ := ecc.NewS256Field(ecc.Secp256k1.Gy)
		_, err := ecc.NewS256Point(x, y)
		if err != nil {
			t.Errorf("NewS256Point: got error %v, expected nil", err)
		}
	})

	t.Run("Verify that the order of G is n", func(t *testing.T) {
		result, _ := ecc.G.ScalarMul(ecc.Secp256k1.N)

		if !result.IsInfinity {
			t.Errorf("ScalarMul: got %v, expected infinite", result)
		}
	})
}
