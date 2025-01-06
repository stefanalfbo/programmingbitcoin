// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"fmt"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
	"github.com/stefanalfbo/programmingbitcoin/encoding/base58"
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

func (p *S256Point) String() string {
	if p.IsInfinity {
		return "S256Point(infinity)"
	}
	x, err := NewS256Field(p.x.number)
	if err != nil {
		return ""
	}
	y, err := NewS256Field(p.y.number)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("S256Point(%s, %s)", x.String(), y.String())
}

// ScalarMul multiplies a point by a scalar
func (p *S256Point) ScalarMul(coefficient *big.Int) (*S256Point, error) {

	c := new(big.Int).Mod(coefficient, Secp256k1.N)

	scalarPoint, err := p.Point.ScalarMul(c)
	if err != nil {
		return nil, err
	}

	return &S256Point{*scalarPoint}, nil
}

func (p *S256Point) Verify(z *big.Int, sig *Signature) (bool, error) {
	s_inv := new(big.Int).Exp(sig.S, new(big.Int).Sub(Secp256k1.N, big.NewInt(2)), Secp256k1.N)
	u := new(big.Int).Mod(new(big.Int).Mul(z, s_inv), Secp256k1.N)
	v := new(big.Int).Mod(new(big.Int).Mul(sig.R, s_inv), Secp256k1.N)

	ug, err := G.ScalarMul(u)
	if err != nil {
		return false, err
	}
	v_point, err := p.ScalarMul(v)
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
func (p *S256Point) SEC() []byte {
	if p.IsInfinity {
		return []byte{0x00}
	}

	return append([]byte{0x04}, append(p.x.number.Bytes(), p.y.number.Bytes()...)...)
}

// Returns the binary representation of the compressed SEC (Standards for Efficient Cryptography) format
func (p *S256Point) SECCompressed() []byte {
	if p.IsInfinity {
		return []byte{0x00}
	}

	if p.y.number.Bit(0) == 0 {
		return append([]byte{0x02}, p.x.number.Bytes()...)
	}

	return append([]byte{0x03}, p.x.number.Bytes()...)
}

// Parse parses a binary representation of the SEC (Standards for Efficient Cryptography) format
func Parse(sec []byte) (*S256Point, error) {
	// if len(sec) == 0 {
	// 	return NewInfinityPoint(), nil
	// }

	if isUncompressed(sec) {
		x := new(big.Int).SetBytes(sec[1:33])
		y := new(big.Int).SetBytes(sec[33:65])

		xField, err := NewS256Field(x)
		if err != nil {
			return nil, err
		}
		yField, err := NewS256Field(y)
		if err != nil {
			return nil, err
		}

		return NewS256Point(xField, yField)
	}

	isEven := sec[0] == 2

	x := new(big.Int).SetBytes(sec[1:])
	xField, err := NewS256Field(x)
	if err != nil {
		return nil, err
	}

	// Right side of the equation y^2 = x^3 + 7
	b, err := NewS256Field(Secp256k1.B)
	if err != nil {
		return nil, err
	}
	x3, err := xField.Pow(big.NewInt(3))
	if err != nil {
		return nil, err
	}
	alpha, err := x3.Add(&b.FieldElement)
	if err != nil {
		return nil, err
	}
	// Solve for left side of the equation
	beta, err := alpha.Sqrt()
	if err != nil {
		return nil, err
	}

	var evenBeta, oddBeta *S256Field
	if beta.number.Bit(0) == 0x00 {
		evenBeta, err = NewS256Field(beta.number)
		if err != nil {
			return nil, err
		}
		oddBeta, err = NewS256Field(new(big.Int).Sub(Secp256k1.Prime, beta.number))
		if err != nil {
			return nil, err
		}
	} else {
		evenBeta, err = NewS256Field(new(big.Int).Sub(Secp256k1.Prime, beta.number))
		if err != nil {
			return nil, err
		}
		oddBeta, err = NewS256Field(beta.number)
		if err != nil {
			return nil, err
		}
	}

	if isEven {
		return NewS256Point(xField, evenBeta)
	}

	return NewS256Point(xField, oddBeta)
}

func isUncompressed(sec []byte) bool {
	return sec[0] == 0x04
}

func (p *S256Point) Hash160(isCompressed bool) []byte {
	if isCompressed {
		return hash.Hash160(p.SECCompressed())
	}

	return hash.Hash160(p.SEC())
}

// Address returns the address of the point
func (p *S256Point) Address(isCompressed bool, isTestnet bool) string {
	hash := p.Hash160(isCompressed)

	if isTestnet {
		return base58.Checksum(append([]byte{0x6f}, hash...))
	}

	return base58.Checksum(append([]byte{0x00}, hash...))
}
