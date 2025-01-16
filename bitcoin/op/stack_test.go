package op_test

import (
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
)

func TestInstruction(t *testing.T) {
	t.Run("Create a too large instruction", func(t *testing.T) {
		data := make([]byte, 521)
		_, err := op.NewInstruction(data)
		if err == nil || err.Error() != "instruction too large" {
			t.Errorf("NewInstruction: expected error, got nil or wrong error")
		}
	})

	t.Run("Create an element", func(t *testing.T) {
		data := []byte{0x01}
		_, err := op.NewInstruction(data)
		if err != nil {
			t.Errorf("NewInstruction: got error %v, expected nil", err)
		}
	})

	t.Run("Equals", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x01})
		element3, _ := op.NewInstruction([]byte{0x02})

		if !element1.Equals(element2) {
			t.Errorf("expected: true, got false")
		}

		if element1.Equals(element3) {
			t.Errorf("expected: false, got true")
		}
	})

	t.Run("Hex", func(t *testing.T) {
		testCases := []struct {
			data     []byte
			expected string
		}{
			{[]byte{0x01}, "01"},
			{[]byte{0x00}, "00"},
			{[]byte{0xff}, "ff"},
			{[]byte{0x12, 0x34}, "1234"},
			{[]byte{0xDE, 0xAD, 0xBE, 0xEF}, "deadbeef"},
		}

		for _, tc := range testCases {
			element, _ := op.NewInstruction(tc.data)
			hex := element.Hex()
			if hex != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected, hex)
			}
		}
	})

	t.Run("IsZero", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x00})
		element2, _ := op.NewInstruction([]byte{0x01})

		if !element1.IsZero() {
			t.Errorf("expected: true, got false")
		}

		if element2.IsZero() {
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
		element, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()

		stack.Push(element)

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}
	})

	t.Run("Pop", func(t *testing.T) {
		element, _ := op.NewInstruction([]byte{0x01})
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
		element, _ := op.NewInstruction([]byte{0x01})
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

	t.Run("PeekN", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		peekedElement, err := stack.PeekN(1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !peekedElement.Equals(element1) {
			t.Fatalf("expected: %v, got: %v", element1, peekedElement)
		}
		if stack.Size() != 2 {
			t.Errorf("expected %v, got %v", 2, stack.Size())
		}
	})
}
