package op_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
)

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
