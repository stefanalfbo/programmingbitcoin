// Package ecc - Elliptic Curve Cryptography
package ecc

import "math/big"

type S256Point struct {
	Point
}

var G *S256Point

func init() {
	x, err := NewS256Field(Secp256k1.Gx)
	if err != nil {
		panic(err)
	}
	y, err := NewS256Field(Secp256k1.Gy)
	if err != nil {
		panic(err)
	}
	G, err = NewS256Point(x, y)
	if err != nil {
		panic(err)
	}
}

func NewS256Point(x, y *S256Field) (*S256Point, error) {
	a, err := NewS256Field(Secp256k1.A)
	if err != nil {
		return nil, err
	}
	b, err := NewS256Field(Secp256k1.B)
	if err != nil {
		return nil, err
	}
	point, err := NewPoint(x.FieldElement, y.FieldElement, a.FieldElement, b.FieldElement)
	if err != nil {
		return nil, err
	}

	return &S256Point{*point}, nil

}

// ScalarMul multiplies a point by a scalar
func (f *S256Point) ScalarMul(coefficient *big.Int) (*S256Point, error) {

	c := new(big.Int).Mod(coefficient, Secp256k1.N)

	p, err := f.Point.ScalarMul(c)
	if err != nil {
		return nil, err
	}

	return &S256Point{*p}, nil
}
