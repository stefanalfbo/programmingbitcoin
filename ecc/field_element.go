// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"errors"
	"fmt"
	"math/big"
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

// Subtract subtracts two field elements
func (f *FieldElement) Subtract(other *FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, errors.New("cannot subtract two numbers in different Fields")
	}

	number := (f.number - other.number + f.prime) % f.prime
	return NewFieldElement(number, f.prime)
}

// Mul multiplies two field elements
func (f *FieldElement) Mul(other *FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, errors.New("cannot multiply two numbers in different Fields")
	}

	number := (f.number * other.number) % f.prime
	return NewFieldElement(number, f.prime)
}

// Pow raises a field element to a power.
func (f *FieldElement) Pow(exponent int) (*FieldElement, error) {
	number := new(big.Int).Exp(
		big.NewInt(int64(f.number)),
		big.NewInt(int64(exponent)),
		big.NewInt(int64(f.prime)))
	return NewFieldElement(int(number.Int64()), f.prime)
}

// Div divides two field elements.
func (f *FieldElement) Div(other *FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, errors.New("cannot divide two numbers in different Fields")
	}

	/*
		User Fermat's Little Theorem to calculate the division.

		number^(prime-1) % prime = 1
		which means:

		1/number == number^(prime-2) % prime
	*/
	pow, err := other.Pow(f.prime - 2)
	if err != nil {
		return nil, err
	}

	return f.Mul(pow)
}
