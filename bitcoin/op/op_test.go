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
