package ecc_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
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

	t.Run("Verify signature", func(t *testing.T) {
		z := new(big.Int).SetBytes([]byte{0xbc, 0x62, 0xd4, 0xb8, 0x0d, 0x9e, 0x36, 0xda, 0x29, 0xc1, 0x6c, 0x5d, 0x4d, 0x9f, 0x11, 0x73, 0x1f, 0x36, 0x05, 0x2c, 0x72, 0x40, 0x1a, 0x76, 0xc2, 0x3c, 0x0f, 0xb5, 0xa9, 0xb7, 0x44, 0x23})
		px, _ := ecc.NewS256Field(new(big.Int).SetBytes([]byte{0x04, 0x51, 0x9f, 0xac, 0x3d, 0x91, 0x0c, 0xa7, 0xe7, 0x13, 0x8f, 0x70, 0x13, 0x70, 0x6f, 0x61, 0x9f, 0xa8, 0xf0, 0x33, 0xe6, 0xec, 0x6e, 0x09, 0x37, 0x0e, 0xa3, 0x8c, 0xee, 0x6a, 0x75, 0x74}))
		py, _ := ecc.NewS256Field(new(big.Int).SetBytes([]byte{0x82, 0xb5, 0x1e, 0xab, 0x8c, 0x27, 0xc6, 0x6e, 0x26, 0xc8, 0x58, 0xa0, 0x79, 0xbc, 0xdf, 0x4f, 0x1a, 0xda, 0x34, 0xce, 0xc4, 0x20, 0xca, 0xfc, 0x7e, 0xac, 0x1a, 0x42, 0x21, 0x6f, 0xb6, 0xc4}))
		point, _ := ecc.NewS256Point(px, py)
		sig := ecc.Signature{
			R: new(big.Int).SetBytes([]byte{0x37, 0x20, 0x6a, 0x06, 0x10, 0x99, 0x5c, 0x58, 0x07, 0x49, 0x99, 0xcb, 0x97, 0x67, 0xb8, 0x7a, 0xf4, 0xc4, 0x97, 0x8d, 0xb6, 0x8c, 0x06, 0xe8, 0xe6, 0xe8, 0x1d, 0x28, 0x20, 0x47, 0xa7, 0xc6}),
			S: new(big.Int).SetBytes([]byte{0x8c, 0xa6, 0x37, 0x59, 0xc1, 0x15, 0x7e, 0xbe, 0xae, 0xc0, 0xd0, 0x3c, 0xec, 0xca, 0x11, 0x9f, 0xc9, 0xa7, 0x5b, 0xf8, 0xe6, 0xd0, 0xfa, 0x65, 0xc8, 0x41, 0xc8, 0xe2, 0x73, 0x8c, 0xda, 0xec}),
		}

		valid, err := point.Verify(z, &sig)
		if err != nil {
			t.Errorf("Verify: got error %v, expected nil", err)
		}
		if !valid {
			t.Errorf("Verify: got %v, expected true", valid)
		}
	})

	t.Run("SEC", func(t *testing.T) {
		coefficient := new(big.Int).Exp(big.NewInt(999), big.NewInt(3), nil)
		uncompressed := "049d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d56fa15cc7f3d38cda98dee2419f415b7513dde1301f8643cd9245aea7f3f911f9"
		point, _ := ecc.G.ScalarMul(coefficient)

		sec := point.SEC()

		if hex.EncodeToString(sec) != uncompressed {
			t.Errorf("SEC: got %v, expected %v", hex.EncodeToString(sec), uncompressed)
		}
	})

	t.Run("SEC compressed", func(t *testing.T) {
		coefficient := new(big.Int).Exp(big.NewInt(999), big.NewInt(3), nil)
		compressed := "039d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d5"
		point, _ := ecc.G.ScalarMul(coefficient)

		secCompressed := point.SECCompressed()

		if hex.EncodeToString(secCompressed) != compressed {
			t.Errorf("SECCompressed: got %v, expected %v", hex.EncodeToString(secCompressed), compressed)
		}
	})

	t.Run("Parse uncompressed SEC format", func(t *testing.T) {
		coefficient := new(big.Int).Exp(big.NewInt(999), big.NewInt(3), nil)
		point, _ := ecc.G.ScalarMul(coefficient)
		sec := point.SEC()

		parsedPoint, err := ecc.Parse(sec)
		if err != nil {
			t.Errorf("Parse: got error %v, expected nil", err)
		}

		if parsedPoint.String() != point.String() {
			t.Errorf("Parse: got %v, expected %v", parsedPoint.String(), point.String())
		}
	})

	t.Run("Parse compressed SEC format", func(t *testing.T) {
		coefficient := new(big.Int).Exp(big.NewInt(999), big.NewInt(3), nil)
		point, _ := ecc.G.ScalarMul(coefficient)
		secCompressed := point.SECCompressed()

		parsedPoint, err := ecc.Parse(secCompressed)
		if err != nil {
			t.Errorf("Parse: got error %v, expected nil", err)
		}

		if parsedPoint.String() != point.String() {
			t.Errorf("Parse: got %v, expected %v", parsedPoint.String(), point.String())
		}
	})
}

func TestVerifyingASignature(t *testing.T) {
	z := new(big.Int).SetBytes([]byte{0xbc, 0x62, 0xd4, 0xb8, 0x0d, 0x9e, 0x36, 0xda, 0x29, 0xc1, 0x6c, 0x5d, 0x4d, 0x9f, 0x11, 0x73, 0x1f, 0x36, 0x05, 0x2c, 0x72, 0x40, 0x1a, 0x76, 0xc2, 0x3c, 0x0f, 0xb5, 0xa9, 0xb7, 0x44, 0x23})
	px, _ := ecc.NewS256Field(new(big.Int).SetBytes([]byte{0x04, 0x51, 0x9f, 0xac, 0x3d, 0x91, 0x0c, 0xa7, 0xe7, 0x13, 0x8f, 0x70, 0x13, 0x70, 0x6f, 0x61, 0x9f, 0xa8, 0xf0, 0x33, 0xe6, 0xec, 0x6e, 0x09, 0x37, 0x0e, 0xa3, 0x8c, 0xee, 0x6a, 0x75, 0x74}))
	py, _ := ecc.NewS256Field(new(big.Int).SetBytes([]byte{0x82, 0xb5, 0x1e, 0xab, 0x8c, 0x27, 0xc6, 0x6e, 0x26, 0xc8, 0x58, 0xa0, 0x79, 0xbc, 0xdf, 0x4f, 0x1a, 0xda, 0x34, 0xce, 0xc4, 0x20, 0xca, 0xfc, 0x7e, 0xac, 0x1a, 0x42, 0x21, 0x6f, 0xb6, 0xc4}))
	point, _ := ecc.NewS256Point(px, py)
	s := new(big.Int).SetBytes([]byte{0x8c, 0xa6, 0x37, 0x59, 0xc1, 0x15, 0x7e, 0xbe, 0xae, 0xc0, 0xd0, 0x3c, 0xec, 0xca, 0x11, 0x9f, 0xc9, 0xa7, 0x5b, 0xf8, 0xe6, 0xd0, 0xfa, 0x65, 0xc8, 0x41, 0xc8, 0xe2, 0x73, 0x8c, 0xda, 0xec})
	s_inv := new(big.Int).Exp(s, new(big.Int).Sub(ecc.Secp256k1.N, big.NewInt(2)), ecc.Secp256k1.N)

	r := new(big.Int).SetBytes([]byte{0x37, 0x20, 0x6a, 0x06, 0x10, 0x99, 0x5c, 0x58, 0x07, 0x49, 0x99, 0xcb, 0x97, 0x67, 0xb8, 0x7a, 0xf4, 0xc4, 0x97, 0x8d, 0xb6, 0x8c, 0x06, 0xe8, 0xe6, 0xe8, 0x1d, 0x28, 0x20, 0x47, 0xa7, 0xc6})

	u := new(big.Int).Mod(new(big.Int).Mul(z, s_inv), ecc.Secp256k1.N)
	v := new(big.Int).Mod(new(big.Int).Mul(r, s_inv), ecc.Secp256k1.N)

	ug, _ := ecc.G.ScalarMul(u)
	vpoint, _ := point.ScalarMul(v)
	result, _ := ug.Add(&vpoint.Point)

	if result.XNum().Cmp(r) != 0 {
		t.Errorf("got %v, expected %v", result.XNum(), r)
	}
}

func TestCreateASignature(t *testing.T) {
	e := new(big.Int).SetBytes(hash.Hash256([]byte("my secret")))
	point, _ := ecc.G.ScalarMul(e)

	expected := "S256Point(028d003eab2e428d11983f3e97c3fa0addf3b42740df0d211795ffb3be2f6c52, 0ae987b9ec6ea159c78cb2a937ed89096fb218d9e7594f02b547526d8cd309e2)"
	if point.String() != expected {
		t.Errorf("got %v, expected %v", point.String(), expected)
	}

	z := new(big.Int).SetBytes(hash.Hash256([]byte("my message")))
	k := big.NewInt(1234567890)
	kG, _ := ecc.G.ScalarMul(k)
	r := kG.XNum()
	k_inv := new(big.Int).Exp(k, new(big.Int).Sub(ecc.Secp256k1.N, big.NewInt(2)), ecc.Secp256k1.N)
	s := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Add(z, new(big.Int).Mul(e, r)), k_inv), ecc.Secp256k1.N)

	rAsHexString := hex.EncodeToString(r.Bytes())
	expectedR := "2b698a0f0a4041b77e63488ad48c23e8e8838dd1fb7520408b121697b782ef22"
	if rAsHexString != expectedR {
		t.Errorf("got %s, expected %s", rAsHexString, expectedR)
	}

	sAsHexString := hex.EncodeToString(s.Bytes())
	expectedS := "bb14e602ef9e3f872e25fad328466b34e6734b7a0fcd58b1eb635447ffae8cb9"
	if sAsHexString != expectedS {
		t.Errorf("got %s, expected %s", sAsHexString, expectedS)
	}
}
