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
