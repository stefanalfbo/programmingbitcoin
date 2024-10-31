package ecc_test

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

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

	t.Run("Pow", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(3, 13)
		b, _ := a.Pow(3)

		expected, _ := ecc.NewFieldElement(1, 13)
		if !b.Equals(expected) {
			t.Errorf("Pow: got %v, expected %v", b, expected)
		}
	})

	t.Run("Pow with negative exponent", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(7, 13)
		b, _ := a.Pow(-3)

		expected, _ := ecc.NewFieldElement(8, 13)
		if !b.Equals(expected) {
			t.Errorf("Pow: got %v, expected %v", b, expected)
		}
	})

	t.Run("Div", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(3, 13)
		b, _ := ecc.NewFieldElement(12, 13)

		c, _ := a.Div(b)

		expected, _ := ecc.NewFieldElement(10, 13)
		if !c.Equals(expected) {
			t.Errorf("Div: got %v, expected %v", c, expected)
		}
	})
}

func TestQuickProperties(t *testing.T) {
	const prime = 257

	makeRandomElementField := func(rnd *rand.Rand, except *ecc.FieldElement) *ecc.FieldElement {
		for {
			f, err := ecc.NewFieldElement(rnd.Intn(prime), prime)
			if err != nil {
				panic(err)
			}

			if except == nil || !f.Equals(except) {
				return f
			}
		}
	}

	generateRandomElementFields := func(length int, except *ecc.FieldElement) func(output []reflect.Value, rnd *rand.Rand) {
		return func(output []reflect.Value, rnd *rand.Rand) {
			for i := 0; i < length; i++ {
				output[i] = reflect.ValueOf(makeRandomElementField(rnd, except))
			}
		}
	}

	t.Run("associativity", func(t *testing.T) {
		generator := generateRandomElementFields(3, nil)
		f := func(a *ecc.FieldElement, b *ecc.FieldElement, c *ecc.FieldElement) bool {
			// (a + b) + c = a + (b + c)
			left1, _ := a.Add(b)
			left, _ := left1.Add(c)

			right1, _ := b.Add(c)
			right, _ := a.Add(right1)

			return left.Equals(right)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})

	t.Run("add commutativity", func(t *testing.T) {
		generator := generateRandomElementFields(2, nil)
		f := func(a *ecc.FieldElement, b *ecc.FieldElement) bool {
			// a + b = b + a
			left, _ := a.Add(b)
			right, _ := b.Add(a)

			return left.Equals(right)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})
}
