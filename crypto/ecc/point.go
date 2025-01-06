// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"errors"
	"fmt"
	"math/big"
)

type Point struct {
	x,
	y,
	a,
	b FieldElement
	IsInfinity bool
}

// NewPoint creates a new point on the elliptic curve
func NewPoint(x, y, a, b FieldElement) (*Point, error) {
	// The canonical form of an elliptic curve is:
	// y^2 = x^3 + ax + b
	left := y.PowUnsafe(big.NewInt(2))
	right := x.PowUnsafe(big.NewInt(3)).AddUnsafe(x.MulUnsafe(&a)).AddUnsafe(&b)

	if !left.Equals(right) {
		return nil, errors.New("point is not on the curve")
	}

	isInfinity := false
	return &Point{x, y, a, b, isInfinity}, nil
}

// NewInfinityPoint creates a new point at infinity
func NewInfinityPoint() *Point {
	return &Point{IsInfinity: true}
}

// String returns the string representation of the point
func (f *Point) String() string {
	if f.IsInfinity {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(%s, %s)_%s_%s", f.x.String(), f.y.String(), f.a.String(), f.b.String())
}

// Equals checks if two points are equal
func (f *Point) Equals(other *Point) bool {
	if f.IsInfinity && other.IsInfinity {
		return true
	} else if f.IsInfinity || other.IsInfinity {
		return false
	}

	return f.x.Equals(&other.x) && f.y.Equals(&other.y) && f.a.Equals(&other.a) && f.b.Equals(&other.b)
}

// Add adds two points
func (f *Point) Add(other *Point) (*Point, error) {
	if !f.IsInfinity && !other.IsInfinity && (!f.a.Equals(&other.a) || !f.b.Equals(&other.b)) {
		return nil, errors.New("cannot add two points which are not on the same curve")
	}

	if f.IsInfinity {
		return other, nil
	}

	if other.IsInfinity {
		return f, nil
	}

	// Additive inverses
	if f.x.Equals(&other.x) && !f.y.Equals(&other.y) {
		return NewInfinityPoint(), nil
	}

	if !f.x.Equals(&other.x) {
		y := other.y.SubtractUnsafe(&f.y)
		x := other.x.SubtractUnsafe(&f.x)
		slope := y.DivUnsafe(x)
		x3 := slope.PowUnsafe(big.NewInt(2)).SubtractUnsafe(&f.x).SubtractUnsafe(&other.x)
		y3 := slope.MulUnsafe(f.x.SubtractUnsafe(x3)).SubtractUnsafe(&f.y)

		return NewPoint(*x3, *y3, f.a, f.b)
	}

	zero := f.x.ScalarMulUnsafe(0) //NewFieldElement(big.NewInt(0), f.x.prime)

	if f.Equals(other) && f.y.Equals(zero) {
		return NewInfinityPoint(), nil
	}

	if f.Equals(other) {
		dividend := f.x.PowUnsafe(big.NewInt(2)).ScalarMulUnsafe(3).AddUnsafe(&f.a)
		divisor := f.y.ScalarMulUnsafe(2)
		slope := dividend.DivUnsafe(divisor)
		x3 := slope.PowUnsafe(big.NewInt(2)).SubtractUnsafe(f.x.ScalarMulUnsafe(2))
		y3 := slope.MulUnsafe(f.x.SubtractUnsafe(x3)).SubtractUnsafe(&f.y)

		return NewPoint(*x3, *y3, f.a, f.b)
	}

	return nil, errors.New("there is no support for adding these points")
}

// ScalarMul multiplies a point by a scalar
func (f *Point) ScalarMul(coefficient *big.Int) (*Point, error) {
	result := NewInfinityPoint()

	for coefficient.Cmp(big.NewInt(0)) > 0 {
		if coefficient.Bit(0) == 1 {
			tmp, err := result.Add(f)
			if err != nil {
				return nil, err
			}
			result = tmp
		}
		tmp, err := f.Add(f)
		if err != nil {
			return nil, err
		}
		f = tmp
		coefficient = new(big.Int).Rsh(coefficient, 1)
	}
	return result, nil
}

func (f *Point) XNum() *big.Int {
	return f.x.number
}
