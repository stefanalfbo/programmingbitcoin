package ecc_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestPoint(t *testing.T) {
	t.Run("NewPoint is valid", func(t *testing.T) {
		_, err := ecc.NewPoint(-1, -1, 5, 7)

		if err != nil {
			t.Errorf("NewPoint: got error %v, expected nil", err)
		}
	})

	t.Run("NewPoint is not valid", func(t *testing.T) {
		_, err := ecc.NewPoint(-1, -2, 5, 7)

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
		p, _ := ecc.NewPoint(-1, -1, 5, 7)

		expected := "Point(-1, -1)_5_7"
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
		p1, _ := ecc.NewPoint(-1, -1, 5, 7)
		p2, _ := ecc.NewPoint(-1, -1, 5, 7)

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
		p1, _ := ecc.NewPoint(-1, -1, 5, 7)
		p2 := ecc.NewInfinityPoint()

		if p1.Equals(p2) {
			t.Errorf("Equals: got %v, expected %v", p1, p2)
		}
	})
}
