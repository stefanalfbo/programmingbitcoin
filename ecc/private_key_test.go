package ecc_test

import (
	"crypto/rand"
	"encoding/hex"
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

	t.Run("Find the uncompressed SEC format for the given private keys", func(t *testing.T) {
		tests := []struct {
			secret       *big.Int
			uncompressed string
		}{
			{
				secret:       big.NewInt(5000),
				uncompressed: "04ffe558e388852f0120e46af2d1b370f85854a8eb0841811ece0e3e03d282d57c315dc72890a4f10a1481c031b03b351b0dc79901ca18a00cf009dbdb157a1d10",
			},
			{
				secret:       new(big.Int).Exp(big.NewInt(2018), big.NewInt(5), nil),
				uncompressed: "04027f3da1918455e03c46f659266a1bb5204e959db7364d2f473bdf8f0a13cc9dff87647fd023c13b4a4994f17691895806e1b40b57f4fd22581a4f46851f3b06",
			},
			{
				secret:       big.NewInt(0xdeadbeef12345),
				uncompressed: "04d90cd625ee87dd38656dd95cf79f65f60f7273b67d3096e68bd81e4f5342691f842efa762fd59961d0e99803c61edba8b3e3f7dc3a341836f97733aebf987121",
			},
		}

		for _, test := range tests {
			privateKey, err := ecc.NewPrivateKey(test.secret)
			if err != nil {
				t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
			}

			uncompressed := hex.EncodeToString(privateKey.SECUncompressed())
			if uncompressed != test.uncompressed {
				t.Errorf("Uncompressed: got %v, expected %v", uncompressed, test.uncompressed)
			}
		}
	})
}
