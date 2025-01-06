// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"fmt"
	"math/big"
)

type Signature struct {
	R, S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("Signature(%x, %x)", s.R, s.S)
}

func NewSignature(r, s *big.Int) *Signature {
	return &Signature{r, s}
}

// DER returns the signature in Distinguished Encoding Rules (DER) format
func (signature *Signature) DER() []byte {
	r := signature.R.Bytes()

	// If the high bit is set in r, prepend a 0x00 byte
	if r[0]&0x80 != 0 {
		r = append([]byte{0x00}, r...)
	}

	result := append([]byte{0x02, byte(len(r))}, r...)

	s := signature.S.Bytes()
	// If the high bit is set in s, prepend a 0x00 byte
	if s[0]&0x80 != 0 {
		s = append([]byte{0x00}, s...)
	}

	result = append(result, 0x02, byte(len(s)))
	result = append(result, s...)

	return append([]byte{0x30, byte(len(result))}, result...)
}

func ParseDER(der []byte) (*Signature, error) {
	if der[0] != 0x30 {
		return nil, fmt.Errorf("bad signature")
	}

	length := int(der[1])
	if length != len(der[2:]) {
		return nil, fmt.Errorf("bad signature length")
	}

	if der[2] != 0x02 {
		return nil, fmt.Errorf("bad signature")
	}

	rLength := int(der[3])
	r := der[4 : 4+rLength]

	if der[4+rLength] != 0x02 {
		return nil, fmt.Errorf("bad signature")
	}

	sLength := int(der[5+rLength])
	s := der[6+rLength:]

	if len(der) != 6+rLength+sLength {
		return nil, fmt.Errorf("signature too long")
	}

	rInt := new(big.Int).SetBytes(r)
	sInt := new(big.Int).SetBytes(s)

	return NewSignature(rInt, sInt), nil
}
