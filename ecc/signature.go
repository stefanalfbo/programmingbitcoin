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
