package ecc_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestFieldElement(t *testing.T) {
	t.Run("NewFieldElement is valid", func(t *testing.T) {
		a, err := ecc.NewFieldElement(7, 13)

		if err != nil {
			t.Errorf("NewFieldElement: got error %v, expected nil", err)
		}
		if a == nil {
			t.Errorf("NewFieldElement: got nil, expected valid FieldElement")
		}
	})

	t.Run("NewFieldElement is not valid", func(t *testing.T) {
		_, err := ecc.NewFieldElement(14, 13)

		if err == nil {
			t.Errorf("NewFieldElement: expected error, got nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(7, 13)

		expected := "FieldElement_13(7)"
		if a.String() != expected {
			t.Errorf("String: got %v, expected %v", a.String(), expected)
		}
	})

	t.Run("not equal", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(7, 13)
		b, _ := ecc.NewFieldElement(6, 14)

		if a.Equals(b) {
			t.Errorf("Equals: got %v, expected %v", a, b)
		}
	})

	t.Run("equal", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(7, 13)

		if !a.Equals(a) {
			t.Errorf("Equals: got %v, expected %v", a, a)
		}
	})

	t.Run("Add", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(7, 13)
		b, _ := ecc.NewFieldElement(12, 13)

		c, _ := a.Add(b)

		expected, _ := ecc.NewFieldElement(6, 13)
		if !c.Equals(expected) {
			t.Errorf("Add: got %v, expected %v", c, expected)
		}
	})

	t.Run("Subtract", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(7, 13)
		b, _ := ecc.NewFieldElement(12, 13)

		c, _ := a.Subtract(b)

		expected, _ := ecc.NewFieldElement(8, 13)
		if !c.Equals(expected) {
			t.Errorf("Subtract: got %v, expected %v", c, expected)
		}
	})

	t.Run("Mul", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(3, 13)
		b, _ := ecc.NewFieldElement(12, 13)

		c, _ := a.Mul(b)

		expected, _ := ecc.NewFieldElement(10, 13)
		if !c.Equals(expected) {
			t.Errorf("Mul: got %v, expected %v", c, expected)
		}
	})
}
