package ecc_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestPoint(t *testing.T) {
	prime := 223
	a, _ := ecc.NewFieldElement(0, prime)
	b, _ := ecc.NewFieldElement(7, prime)
	x1, _ := ecc.NewFieldElement(192, prime)
	y1, _ := ecc.NewFieldElement(105, prime)
	x2, _ := ecc.NewFieldElement(17, prime)
	y2, _ := ecc.NewFieldElement(56, prime)

	t.Run("NewPoint is valid", func(t *testing.T) {
		_, err := ecc.NewPoint(*x1, *y1, *a, *b)

		if err != nil {
			t.Errorf("NewPoint: got error %v, expected nil", err)
		}
	})

	t.Run("NewPoint is not valid", func(t *testing.T) {
		_, err := ecc.NewPoint(*x1, *y1, *a, *a)

		if err == nil {
			t.Errorf("NewPoint: expected error, got nil")
		}
	})

	t.Run("NewPoint is infinite", func(t *testing.T) {
		p := ecc.NewInfinityPoint()

		if !p.IsInfinity {
			t.Errorf("NewPoint: got %v, expected infinite", p)
		}
	})

	t.Run("String", func(t *testing.T) {
		p, _ := ecc.NewPoint(*x1, *y1, *a, *b)

		expected := "Point({192 223}, {105 223})_{0 223}_{7 223}"
		if p.String() != expected {
			t.Errorf("String: got %v, expected %v", p.String(), expected)
		}
	})

	t.Run("String infinite", func(t *testing.T) {
		p := ecc.NewInfinityPoint()

		expected := "Point(infinity)"
		if p.String() != expected {
			t.Errorf("String: got %v, expected %v", p.String(), expected)
		}
	})

	t.Run("Equals", func(t *testing.T) {
		p1, _ := ecc.NewPoint(*x1, *y1, *a, *b)
		p2, _ := ecc.NewPoint(*x1, *y1, *a, *b)

		if !p1.Equals(p2) {
			t.Errorf("Equals: got %v, expected %v", p1, p2)
		}
	})

	t.Run("Equals infinity point", func(t *testing.T) {
		p1 := ecc.NewInfinityPoint()
		p2 := ecc.NewInfinityPoint()

		if !p1.Equals(p2) {
			t.Errorf("Equals: got %v, expected %v", p1, p2)
		}
	})

	t.Run("not equal", func(t *testing.T) {
		p1, _ := ecc.NewPoint(*x1, *y1, *a, *b)
		p2 := ecc.NewInfinityPoint()

		if p1.Equals(p2) {
			t.Errorf("Equals: got %v, expected %v", p1, p2)
		}
	})

	t.Run("add points on the curve", func(t *testing.T) {
		p1, _ := ecc.NewPoint(*x1, *y1, *a, *b)
		p2, _ := ecc.NewPoint(*x2, *y2, *a, *b)
		x3, _ := ecc.NewFieldElement(170, prime)
		y3, _ := ecc.NewFieldElement(142, prime)
		expected, _ := ecc.NewPoint(*x3, *y3, *a, *b)

		result, err := p1.Add(p2)
		if err != nil {
			t.Errorf("Add: got error %v, expected nil", err)
		}
		if !result.Equals(expected) {
			t.Errorf("Add: got %v, expected %v", result, expected)
		}
	})

	t.Run("add two points not on the same curve", func(t *testing.T) {
		p1, _ := ecc.NewPoint(*x1, *y1, *a, *b)
		a1, _ := ecc.NewFieldElement(5, prime)
		b1, _ := ecc.NewFieldElement(7, prime)
		x2, _ := ecc.NewFieldElement(18, prime)
		y2, _ := ecc.NewFieldElement(77, prime)
		p2, _ := ecc.NewPoint(*x2, *y2, *a1, *b1)

		_, err := p1.Add(p2)
		if err == nil {
			t.Errorf("Add: expected error, got nil")
		}
	})

	t.Run("add point to infinity", func(t *testing.T) {
		p1 := ecc.NewInfinityPoint()
		p2, _ := ecc.NewPoint(*x1, *y1, *a, *b)

		result, err := p1.Add(p2)
		if err != nil {
			t.Errorf("Add: got error %v, expected nil", err)
		}
		if !result.Equals(p2) {
			t.Errorf("Add: got %v, expected %v", result, p2)
		}
	})
}
