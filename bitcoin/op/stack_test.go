package op_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
)

func TestElement(t *testing.T) {
	t.Run("Create a too small element", func(t *testing.T) {
		_, err := op.NewElement([]byte{})
		if err == nil || err.Error() != "element too small" {
			t.Errorf("NewElement: expected error, got nil")
		}
	})

	t.Run("Create a too large element", func(t *testing.T) {
		data := make([]byte, 521)
		_, err := op.NewElement(data)
		if err == nil || err.Error() != "element too large" {
			t.Errorf("NewElement: expected error, got nil")
		}
	})

	t.Run("Create an element", func(t *testing.T) {
		data := []byte{0x01}
		_, err := op.NewElement(data)
		if err != nil {
			t.Errorf("NewElement: got error %v, expected nil", err)
		}
	})

	t.Run("Equals", func(t *testing.T) {
		element1, _ := op.NewElement([]byte{0x01})
		element2, _ := op.NewElement([]byte{0x01})
		element3, _ := op.NewElement([]byte{0x02})

		if !element1.Equals(element2) {
			t.Errorf("expected: true, got false")
		}

		if element1.Equals(element3) {
			t.Errorf("expected: false, got true")
		}
	})
}

func TestStack(t *testing.T) {
	t.Run("Create new stack", func(t *testing.T) {
		stack := op.NewStack()

		if stack.Size() != 0 {
			t.Errorf("expected: %v, got: %v", 0, stack.Size())
		}
	})
	t.Run("Push", func(t *testing.T) {
		element, _ := op.NewElement([]byte{0x01})
		stack := op.NewStack()

		stack.Push(element)

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}
	})

	t.Run("Pop", func(t *testing.T) {
		element, _ := op.NewElement([]byte{0x01})
		stack := op.NewStack()
		stack.Push(element)

		poppedElement, err := stack.Pop()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !poppedElement.Equals(element) {
			t.Fatalf("expected: %v, got: %v", element, poppedElement)
		}
		if stack.Size() != 0 {
			t.Errorf("expected %v, got %v", 0, stack.Size())
		}
	})

	t.Run("Peek", func(t *testing.T) {
		element, _ := op.NewElement([]byte{0x01})
		stack := op.NewStack()
		stack.Push(element)

		peekedElement, err := stack.Peek()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !peekedElement.Equals(element) {
			t.Fatalf("expected: %v, got: %v", element, peekedElement)
		}
		if stack.Size() != 1 {
			t.Errorf("expected %v, got %v", 1, stack.Size())
		}
	})

}
