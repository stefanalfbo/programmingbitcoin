// Package ecc - Elliptic Curve Cryptography
package ecc

import (
	"errors"
	"fmt"
)

type Point struct {
	x,
	y,
	a,
	b int
}

// NewPoint creates a new point on the elliptic curve
func NewPoint(x, y, a, b int) (*Point, error) {
	// The canonical form of an elliptic curve is:
	// y^2 = x^3 + ax + b
	if y*y != x*x*x+a*x+b {
		return nil, errors.New("point is not on the curve")
	}

	return &Point{x, y, a, b}, nil
}

// String returns the string representation of the point
func (f *Point) String() string {
	return fmt.Sprintf("Point(%d, %d)_%d_%d", f.x, f.y, f.a, f.b)
}

// Equals checks if two points are equal
func (f *Point) Equals(other *Point) bool {
	return f.x == other.x && f.y == other.y && f.a == other.a && f.b == other.b
}
