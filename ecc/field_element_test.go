package ecc_test

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

func TestFieldElement(t *testing.T) {
	t.Run("NewFieldElement is valid", func(t *testing.T) {
		a, err := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		if err != nil {
			t.Errorf("NewFieldElement: got error %v, expected nil", err)
		}
		if a == nil {
			t.Errorf("NewFieldElement: got nil, expected valid FieldElement")
		}
	})

	t.Run("NewFieldElement is not valid", func(t *testing.T) {
		_, err := ecc.NewFieldElement(big.NewInt(14), big.NewInt(13))

		if err == nil {
			t.Errorf("NewFieldElement: expected error, got nil")
		}
	})

	t.Run("String", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))

		expected := "FieldElement_13(7)"
		if a.String() != expected {
			t.Errorf("String: got %v, expected %v", a.String(), expected)
		}
	})

	t.Run("not equal", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(6), big.NewInt(14))

		if a.Equals(b) {
			t.Errorf("Equals: got %v, expected %v", a, b)
		}
	})

	t.Run("equal", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))

		if !a.Equals(a) {
			t.Errorf("Equals: got %v, expected %v", a, a)
		}
	})

	t.Run("Add", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c, _ := a.Add(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(6), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("Add: got %v, expected %v", c, expected)
		}
	})

	t.Run("AddUnsafe", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c := a.AddUnsafe(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(6), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("AddUnsafe: got %v, expected %v", c, expected)
		}
	})

	t.Run("Subtract", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c, _ := a.Subtract(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(8), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("Subtract: got %v, expected %v", c, expected)
		}
	})

	t.Run("SubtractUnsafe", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c := a.SubtractUnsafe(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(8), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("SubtractUnsafe: got %v, expected %v", c, expected)
		}
	})

	t.Run("Mul", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c, _ := a.Mul(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(10), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("Mul: got %v, expected %v", c, expected)
		}
	})

	t.Run("MulUnsafe", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c := a.MulUnsafe(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(10), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("MulUnsafe: got %v, expected %v", c, expected)
		}
	})

	t.Run("ScalarMul", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b, _ := a.ScalarMul(3)

		expected, _ := ecc.NewFieldElement(big.NewInt(9), big.NewInt(13))
		if !b.Equals(expected) {
			t.Errorf("ScalarMul: got %v, expected %v", b, expected)
		}
	})

	t.Run("ScalarMulUnsafe", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b := a.ScalarMulUnsafe(3)

		expected, _ := ecc.NewFieldElement(big.NewInt(9), big.NewInt(13))
		if !b.Equals(expected) {
			t.Errorf("ScalarMulUnsafe: got %v, expected %v", b, expected)
		}
	})

	t.Run("Pow", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b, _ := a.Pow(big.NewInt(3))

		expected, _ := ecc.NewFieldElement(big.NewInt(1), big.NewInt(13))
		if !b.Equals(expected) {
			t.Errorf("Pow: got %v, expected %v", b, expected)
		}
	})

	t.Run("PowUnsafe", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b := a.PowUnsafe(big.NewInt(3))

		expected, _ := ecc.NewFieldElement(big.NewInt(1), big.NewInt(13))
		if !b.Equals(expected) {
			t.Errorf("Pow: got %v, expected %v", b, expected)
		}
	})

	t.Run("Pow with negative exponent", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(7), big.NewInt(13))
		b, _ := a.Pow(big.NewInt(-3))

		expected, _ := ecc.NewFieldElement(big.NewInt(8), big.NewInt(13))
		if !b.Equals(expected) {
			t.Errorf("Pow: got %v, expected %v", b, expected)
		}
	})

	t.Run("Div", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c, _ := a.Div(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(10), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("Div: got %v, expected %v", c, expected)
		}
	})

	t.Run("DivUnsafe", func(t *testing.T) {
		a, _ := ecc.NewFieldElement(big.NewInt(3), big.NewInt(13))
		b, _ := ecc.NewFieldElement(big.NewInt(12), big.NewInt(13))

		c := a.DivUnsafe(b)

		expected, _ := ecc.NewFieldElement(big.NewInt(10), big.NewInt(13))
		if !c.Equals(expected) {
			t.Errorf("DivUnsafe: got %v, expected %v", c, expected)
		}
	})
}

func TestQuickProperties(t *testing.T) {
	const prime = 257

	makeRandomElementField := func(rnd *rand.Rand, except *ecc.FieldElement) *ecc.FieldElement {
		for {
			// Generate a random number between 1 and prime-1
			number := rnd.Intn(prime-1) + 1
			f, err := ecc.NewFieldElement(big.NewInt(int64(number)), big.NewInt(prime))
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

	t.Run("multiplication commutativity", func(t *testing.T) {
		generator := generateRandomElementFields(2, nil)
		f := func(a *ecc.FieldElement, b *ecc.FieldElement) bool {
			// a * b = b * a
			left, _ := a.Mul(b)
			right, _ := b.Mul(a)

			return left.Equals(right)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})

	t.Run("add identity", func(t *testing.T) {
		generator := generateRandomElementFields(1, nil)
		f := func(a *ecc.FieldElement) bool {
			// a + 0 = a
			zero, _ := ecc.NewFieldElement(big.NewInt(0), big.NewInt(prime))

			left, _ := a.Add(zero)

			return left.Equals(a)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})

	t.Run("multiplication identity", func(t *testing.T) {
		generator := generateRandomElementFields(1, nil)
		f := func(a *ecc.FieldElement) bool {
			// a * 1 = a
			one, _ := ecc.NewFieldElement(big.NewInt(1), big.NewInt(prime))

			left, _ := a.Mul(one)

			return left.Equals(a)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})

	t.Run("multiplicative inverse", func(t *testing.T) {
		generator := generateRandomElementFields(1, nil)
		f := func(a *ecc.FieldElement) bool {
			// a * a^-1 = 1
			one, _ := ecc.NewFieldElement(big.NewInt(1), big.NewInt(prime))

			inverse, err := a.Pow(big.NewInt(prime - 2))
			if err != nil {
				return false
			}

			left, _ := a.Mul(inverse)

			return left.Equals(one)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})

	t.Run("distributive property", func(t *testing.T) {
		generator := generateRandomElementFields(3, nil)
		f := func(a *ecc.FieldElement, b *ecc.FieldElement, c *ecc.FieldElement) bool {
			// a * (b + c) = a * b + a * c
			left1, _ := b.Add(c)
			left, _ := a.Mul(left1)

			right1, _ := a.Mul(b)
			right2, _ := a.Mul(c)
			right, _ := right1.Add(right2)

			return left.Equals(right)
		}
		config := quick.Config{Values: generator}
		if err := quick.Check(f, &config); err != nil {
			t.Error(err)
		}
	})
}
