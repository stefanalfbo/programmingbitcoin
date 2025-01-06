// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/encoding/base58"
)

type PrivateKey struct {
	secret *big.Int
	// Public key
	point *S256Point
}

func NewPrivateKey(secret *big.Int) (*PrivateKey, error) {
	point, err := G.ScalarMul(secret)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{secret, point}, nil
}

func (pk *PrivateKey) Hex() string {
	return fmt.Sprintf("%064x", pk.secret)
}

func (pk *PrivateKey) Sign(z *big.Int) (*Signature, error) {
	// k, err := rand.Int(rand.Reader, Secp256k1.N)
	k := deterministicK(z, pk.secret)
	// if err != nil {
	// 	return nil, err
	// }

	kG, err := G.ScalarMul(k)
	if err != nil {
		return nil, err
	}
	r := kG.XNum()

	kInverse := new(big.Int).Exp(k, new(big.Int).Sub(Secp256k1.N, big.NewInt(2)), Secp256k1.N)
	s := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Add(z, new(big.Int).Mul(pk.secret, r)), kInverse), Secp256k1.N)

	if s.Cmp(new(big.Int).Div(Secp256k1.N, big.NewInt(2))) > 0 {
		s = new(big.Int).Sub(Secp256k1.N, s)
	}

	return NewSignature(r, s), nil
}

func deterministicK(z, secret *big.Int) *big.Int {
	k := make([]byte, 32)
	v := bytes.Repeat([]byte{0x01}, 32)

	if z.Cmp(Secp256k1.N) > 0 {
		z = new(big.Int).Sub(z, Secp256k1.N)
	}

	zBytes := z.FillBytes(make([]byte, 32))
	secretBytes := secret.FillBytes(make([]byte, 32))

	h := func(key, data []byte) []byte {
		mac := hmac.New(sha256.New, key)
		mac.Write(data)
		return mac.Sum(nil)
	}

	k = h(k, append(append(v, 0x00), append(secretBytes, zBytes...)...))
	v = h(k, v)
	k = h(k, append(append(v, 0x01), append(secretBytes, zBytes...)...))
	v = h(k, v)

	// Generate candidate values until a valid one is found
	for {
		v = h(k, v)
		candidate := new(big.Int).SetBytes(v)
		if candidate.Cmp(big.NewInt(1)) >= 0 && candidate.Cmp(Secp256k1.N) < 0 {
			return candidate
		}
		k = h(k, append(v, 0x00))
		v = h(k, v)
	}
}

func (pk *PrivateKey) Verify(z *big.Int, signature *Signature) (bool, error) {
	result, err := pk.point.Verify(z, signature)
	if err != nil {
		return false, err
	}

	return result, nil
}

// Returns the binary representation of the SEC (Standards for Efficient Cryptography) format
func (pk *PrivateKey) SECUncompressed() []byte {
	return pk.point.SEC()
}

// Returns the binary representation of the compressed SEC (Standards for Efficient Cryptography) format
func (pk *PrivateKey) SECCompressed() []byte {
	return pk.point.SECCompressed()
}

func (pk *PrivateKey) Address(isCompressed, isTestnet bool) string {
	return pk.point.Address(isCompressed, isTestnet)
}

// WIF (Wallet Import Format) is a way to encode the private key to make it easier to copy
func (pk *PrivateKey) WIF(isCompressed, isTestnet bool) string {
	s := pk.secret.FillBytes(make([]byte, 32))

	if isTestnet {
		s = append([]byte{0xef}, s...)
	} else {
		s = append([]byte{0x80}, s...)
	}

	if isCompressed {
		s = append(s, 0x01)
	}

	return base58.Checksum(s)
}
