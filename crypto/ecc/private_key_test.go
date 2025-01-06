package ecc_test

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
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

	t.Run("Find the compressed SEC format for the given private keys", func(t *testing.T) {
		tests := []struct {
			secret     *big.Int
			compressed string
		}{
			{
				secret:     big.NewInt(5001),
				compressed: "0357a4f368868a8a6d572991e484e664810ff14c05c0fa023275251151fe0e53d1",
			},
			{
				secret:     new(big.Int).Exp(big.NewInt(2019), big.NewInt(5), nil),
				compressed: "02933ec2d2b111b92737ec12f1c5d20f3233a0ad21cd8b36d0bca7a0cfa5cb8701",
			},
			{
				secret:     big.NewInt(0xdeadbeef54321),
				compressed: "0296be5b1292f6c856b3c5654e886fc13511462059089cdf9c479623bfcbe77690",
			},
		}

		for _, test := range tests {
			privateKey, err := ecc.NewPrivateKey(test.secret)
			if err != nil {
				t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
			}

			compressed := hex.EncodeToString(privateKey.SECCompressed())
			if compressed != test.compressed {
				t.Errorf("Compressed: got %v, expected %v", compressed, test.compressed)
			}
		}
	})

	t.Run("Address", func(t *testing.T) {
		tests := []struct {
			secret       *big.Int
			isCompressed bool
			isTestnet    bool
			address      string
		}{
			{
				secret:       big.NewInt(5002),
				isCompressed: false,
				isTestnet:    true,
				address:      "mmTPbXQFxboEtNRkwfh6K51jvdtHLxGeMA",
			},
			{
				secret:       new(big.Int).Exp(big.NewInt(2020), big.NewInt(5), nil),
				isCompressed: true,
				isTestnet:    true,
				address:      "mopVkxp8UhXqRYbCYJsbeE1h1fiF64jcoH",
			},
			{
				secret:       big.NewInt(0x12345deadbeef),
				isCompressed: true,
				isTestnet:    false,
				address:      "1F1Pn2y6pDb68E5nYJJeba4TLg2U7B6KF1",
			},
		}

		for _, test := range tests {
			privateKey, err := ecc.NewPrivateKey(test.secret)
			if err != nil {
				t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
			}

			address := privateKey.Address(test.isCompressed, test.isTestnet)
			if address != test.address {
				t.Errorf("Address: got %v, expected %v", address, test.address)
			}
		}
	})

	t.Run("WIF", func(t *testing.T) {
		tests := []struct {
			secret       *big.Int
			isCompressed bool
			isTestnet    bool
			wif          string
		}{
			{
				secret:       big.NewInt(5003),
				isCompressed: true,
				isTestnet:    true,
				wif:          "cMahea7zqjxrtgAbB7LSGbcQUr1uX1ojuat9jZodMN8rFTv2sfUK",
			},
			{
				secret:       new(big.Int).Exp(big.NewInt(2021), big.NewInt(5), nil),
				isCompressed: false,
				isTestnet:    true,
				wif:          "91avARGdfge8E4tZfYLoxeJ5sGBdNJQH4kvjpWAxgzczjbCwxic",
			},
			{
				secret:       big.NewInt(0x54321deadbeef),
				isCompressed: true,
				isTestnet:    false,
				wif:          "KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgiuQJv1h8Ytr2S53a",
			},
		}

		for _, test := range tests {
			privateKey, err := ecc.NewPrivateKey(test.secret)
			if err != nil {
				t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
			}

			wif := privateKey.WIF(test.isCompressed, test.isTestnet)
			if wif != test.wif {
				t.Errorf("WIF: got %v, expected %v", wif, test.wif)
			}
		}
	})
}

func TestCreateTestnetAddress(t *testing.T) {
	passphrase := "stefan@alfbo.se my secret"
	secret := endian.LittleEndianToBigInt(ecc.Hash256(passphrase).Bytes())
	privateKey, err := ecc.NewPrivateKey(secret)
	if err != nil {
		t.Fatalf("NewPrivateKey: got error %v, expected nil", err)
	}

	address := privateKey.Address(true, true)
	expected := "mkoxbQEyJWRkUoAZVeF79gdMvwqwpCR7mE"
	if address != expected {
		t.Errorf("Address: got %v, expected %v", address, expected)
	}
}
