// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"fmt"
	"math/big"
)

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

func (f *S256Point) String() string {
	if f.IsInfinity {
		return "S256Point(infinity)"
	}
	x, err := NewS256Field(f.x.number)
	if err != nil {
		return ""
	}
	y, err := NewS256Field(f.y.number)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("S256Point(%s, %s)", x.String(), y.String())
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

func (f *S256Point) Verify(z *big.Int, sig *Signature) (bool, error) {
	s_inv := new(big.Int).Exp(sig.S, new(big.Int).Sub(Secp256k1.N, big.NewInt(2)), Secp256k1.N)
	u := new(big.Int).Mod(new(big.Int).Mul(z, s_inv), Secp256k1.N)
	v := new(big.Int).Mod(new(big.Int).Mul(sig.R, s_inv), Secp256k1.N)

	ug, err := G.ScalarMul(u)
	if err != nil {
		return false, err
	}
	v_point, err := f.ScalarMul(v)
	if err != nil {
		return false, err
	}

	total, err := ug.Add(&v_point.Point)
	if err != nil {
		return false, err
	}

	return total.XNum().Cmp(sig.R) == 0, nil
}

// Returns the binary representation of the SEC (Standards for Efficient Cryptography) format
func (f *S256Point) SEC() []byte {
	if f.IsInfinity {
		return []byte{0x00}
	}

	return append([]byte{0x04}, append(f.x.number.Bytes(), f.y.number.Bytes()...)...)
}

// Returns the binary representation of the compressed SEC (Standards for Efficient Cryptography) format
func (f *S256Point) SECCompressed() []byte {
	if f.IsInfinity {
		return []byte{0x00}
	}

	if f.y.number.Bit(0) == 0 {
		return append([]byte{0x02}, f.x.number.Bytes()...)
	}

	return append([]byte{0x03}, f.x.number.Bytes()...)
}
