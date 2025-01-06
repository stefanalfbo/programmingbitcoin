// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"fmt"
	"math/big"
)

type S256Field struct {
	FieldElement
}

func (s *S256Field) String() string {
	return fmt.Sprintf("%064x", s.number)
}

func NewS256Field(number *big.Int) (*S256Field, error) {
	field, err := NewFieldElement(number, Secp256k1.Prime)
	if err != nil {
		return nil, err
	}

	return &S256Field{*field}, nil
}
