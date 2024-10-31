// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"errors"
	"fmt"
)

// FieldElement represents a field element in a prime field
type FieldElement struct {
	number int
	prime  int
}

// NewFieldElement creates a new field element
func NewFieldElement(number, prime int) (*FieldElement, error) {
	if number >= prime || number < 0 {
		return nil, errors.New("number not in field range")
	}

	return &FieldElement{number, prime}, nil
}

// String returns the string representation of the field element
func (f *FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", f.prime, f.number)
}

// Equals checks if two field elements are equal
func (f *FieldElement) Equals(other *FieldElement) bool {
	return f.number == other.number && f.prime == other.prime
}

// Add adds two field elements
func (f *FieldElement) Add(other *FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, errors.New("cannot add two numbers in different Fields")
	}

	number := (f.number + other.number) % f.prime
	return NewFieldElement(number, f.prime)
}
