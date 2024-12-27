package ecc_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestPrivateKey(t *testing.T) {

	t.Run("Hex", func(t *testing.T) {
		secret := big.NewInt(12345)
		privateKey, err := ecc.NewPrivateKey(secret)
		if err != nil {
			t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
		}

		expected := "0000000000000000000000000000000000000000000000000000000000003039"
		if privateKey.Hex() != expected {
			t.Errorf("Hex: got %v, expected %v", privateKey.Hex(), expected)
		}
	})

	t.Run("Sign", func(t *testing.T) {
		secret, err := rand.Int(rand.Reader, ecc.Secp256k1.N)
		if err != nil {
			t.Fatalf("rand.Int: got error %v, expected nil", err)
		}
		privateKey, err := ecc.NewPrivateKey(secret)
		if err != nil {
			t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
		}

		z, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 256))
		if err != nil {
			t.Fatalf("rand.Int: got error %v, expected nil", err)
		}

		signature, err := privateKey.Sign(z)
		if err != nil {
			t.Fatalf("Sign: got error %v, expected nil", err)
		}

		valid, err := privateKey.Verify(z, signature)
		if err != nil {
			t.Fatalf("Verify: got error %v, expected nil", err)
		}

		if !valid {
			t.Errorf("Verify: got %v, expected true", valid)
		}
	})
}
