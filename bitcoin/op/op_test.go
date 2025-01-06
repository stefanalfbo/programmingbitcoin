package op_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
)

func TestOP0(t *testing.T) {
	t.Run("OP 0", func(t *testing.T) {
		stack := op.NewStack()
		stack, err := op.OP0(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		element, _ := stack.Pop()
		if element.Hex() != "" {
			t.Errorf("expected: %v, got: %v", "", element.Hex())
		}
	})
}

func TestOP1(t *testing.T) {
	t.Run("OP 1", func(t *testing.T) {
		stack := op.NewStack()
		stack, err := op.OP1(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		element, _ := stack.Pop()
		if element.Hex() != "01" {
			t.Errorf("expected: %v, got: %v", "01", element.Hex())
		}
	})
}

func TestOP16(t *testing.T) {
	t.Run("OP 16", func(t *testing.T) {
		stack := op.NewStack()
		stack, err := op.OP16(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		element, _ := stack.Pop()
		if element.Hex() != "10" {
			t.Errorf("expected: %v, got: %v", "01", element.Hex())
		}
	})
}

func TestDUP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.DUP(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Duplicate element", func(t *testing.T) {
		element, _ := op.NewElement([]byte{0x01})
		stack := op.NewStack()
		stack.Push(element)

		stack, err := op.DUP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 2 {
			t.Errorf("expected: %v, got: %v", 2, stack.Size())
		}
	})
}

func TestADD(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.ADD(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Add elements", func(t *testing.T) {
		element1, _ := op.NewElement([]byte{0x01})
		element2, _ := op.NewElement([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.ADD(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		element, _ := stack.Pop()
		if element.Hex() != "03" {
			t.Errorf("expected: %v, got: %v", "03", element.Hex())
		}
	})
}

func TestHASH160(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.HASH160(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Hash element", func(t *testing.T) {
		expected := "d7d5ee7824ff93f94c3055af9382c86c68b5ca92"
		element, _ := op.NewElement([]byte("hello world"))
		stack := op.NewStack()
		stack.Push(element)

		stack, err := op.HASH160(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		hashedElement, _ := stack.Pop()
		if expected != hashedElement.Hex() {
			t.Errorf("expected: %v, got: %v", expected, hashedElement.Hex())
		}
	})

}
func TestHASH256(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.HASH256(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Hash element", func(t *testing.T) {
		element, _ := op.NewElement([]byte{0x01})
		stack := op.NewStack()
		stack.Push(element)

		stack, err := op.HASH256(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}
	})
}

func TestCHECKSIG(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.CHECKSIG(stack, nil)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Check signature", func(t *testing.T) {
		z, _ := new(big.Int).SetString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 16)
		secBytes, _ := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
		sigBytes, _ := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601")

		sec, _ := op.NewElement(secBytes)
		sig, _ := op.NewElement(sigBytes)

		stack := op.NewStack()
		stack.Push(sig)
		stack.Push(sec)

		stack, err := op.CHECKSIG(stack, z)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		element, _ := stack.Pop()
		if element.Hex() != "01" {
			t.Errorf("expected: %v, got: %v", "01", element.Hex())
		}
	})
}