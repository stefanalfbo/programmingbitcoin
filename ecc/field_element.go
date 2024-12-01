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

// AddUnsafe adds two field elements without error checking
func (f *FieldElement) AddUnsafe(other *FieldElement) *FieldElement {
	sum, err := f.Add(other)
	if err != nil {
		panic(err)
	}

	return sum
}

// Subtract subtracts two field elements
func (f *FieldElement) Subtract(other *FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, errors.New("cannot subtract two numbers in different Fields")
	}

	number := (f.number - other.number + f.prime) % f.prime
	return NewFieldElement(number, f.prime)
}

// SubtractUnsafe subtracts two field elements without error checking
func (f *FieldElement) SubtractUnsafe(other *FieldElement) *FieldElement {
	difference, err := f.Subtract(other)
	if err != nil {
		panic(err)
	}

	return difference
}

// Mul multiplies two field elements
func (f *FieldElement) Mul(other *FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, errors.New("cannot multiply two numbers in different Fields")
	}

	number := (f.number * other.number) % f.prime
	return NewFieldElement(number, f.prime)
}

// MulUnsafe multiplies two field elements without error checking
func (f *FieldElement) MulUnsafe(other *FieldElement) *FieldElement {
	product, err := f.Mul(other)
	if err != nil {
		panic(err)
	}

	return product
}

// Scalar multiplication multiplies a field element by a scalar
func (f *FieldElement) ScalarMul(scalar int) (*FieldElement, error) {
	number := (f.number * scalar) % f.prime
	return NewFieldElement(number, f.prime)
}

// ScalarMulUnsafe multiplies a field element by a scalar without error checking
func (f *FieldElement) ScalarMulUnsafe(scalar int) *FieldElement {
	product, err := f.ScalarMul(scalar)
	if err != nil {
		panic(err)
	}

	return product
}

// Pow raises a field element to a power.
func (f *FieldElement) Pow(exponent int) (*FieldElement, error) {
	number := new(big.Int).Exp(
		big.NewInt(int64(f.number)),
		big.NewInt(int64(exponent)),
		big.NewInt(int64(f.prime)))
	return NewFieldElement(int(number.Int64()), f.prime)
}

// PowUnsafe raises a field element to a power without error checking.
func (f *FieldElement) PowUnsafe(exponent int) *FieldElement {
	pow, err := f.Pow(exponent)
	if err != nil {
		panic(err)
	}

	return pow
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

// DivUnsafe divides two field elements without error checking.
func (f *FieldElement) DivUnsafe(other *FieldElement) *FieldElement {
	quotient, err := f.Div(other)
	if err != nil {
		panic(err)
	}

	return quotient
}
