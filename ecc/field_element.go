// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"errors"
	"fmt"
	"math/big"
)

// FieldElement represents a field element in a prime field
type FieldElement struct {
	number *big.Int
	prime  *big.Int
}

// NewFieldElement creates a new field element
func NewFieldElement(number, prime *big.Int) (*FieldElement, error) {
	if number.Cmp(prime) >= 0 || number.Cmp(big.NewInt(0)) < 0 {
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
	return f.number.Cmp(other.number) == 0 && f.prime.Cmp(other.prime) == 0
}

// Add adds two field elements
func (f *FieldElement) Add(other *FieldElement) (*FieldElement, error) {
	if f.prime.Cmp(other.prime) != 0 {
		return nil, errors.New("cannot add two numbers in different Fields")
	}

	number := new(big.Int).Add(f.number, other.number)
	number.Mod(number, f.prime)
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
	if f.prime.Cmp(other.prime) != 0 {
		return nil, errors.New("cannot subtract two numbers in different Fields")
	}

	number := new(big.Int).Sub(f.number, other.number)
	number.Add(number, f.prime)
	number.Mod(number, f.prime)
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
	if f.prime.Cmp(other.prime) != 0 {
		return nil, errors.New("cannot multiply two numbers in different Fields")
	}

	number := new(big.Int).Mul(f.number, other.number)
	number.Mod(number, f.prime)
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
	scalarBigInt := big.NewInt(int64(scalar))
	number := new(big.Int).Mul(f.number, scalarBigInt)
	number.Mod(number, f.prime)
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
	exponentBigInt := big.NewInt(int64(exponent))
	n := new(big.Int).Mod(exponentBigInt, new(big.Int).Sub(f.prime, big.NewInt(1)))
	number := new(big.Int).Exp(
		f.number,
		n,
		f.prime)
	return NewFieldElement(number, f.prime)
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
	if f.prime.Cmp(other.prime) != 0 {
		return nil, errors.New("cannot divide two numbers in different Fields")
	}

	/*
		User Fermat's Little Theorem to calculate the division.

		number^(prime-1) % prime = 1
		which means:

		1/number == number^(prime-2) % prime
	*/
	two := big.NewInt(2)
	primeMinusTwo := new(big.Int).Sub(f.prime, two)
	pow, err := other.Pow(int(primeMinusTwo.Int64()))
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
