package op_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
)

func TestOP0_OP16(t *testing.T) {
	t.Run("OP 0", func(t *testing.T) {
		expected := ""
		stack := op.NewStack()
		stack, err := op.OP0(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ := stack.Pop()
		if instruction.Hex() != expected {
			t.Errorf("expected: %v, got: %v", expected, instruction.Hex())
		}
	})

	t.Run("OP 1NEGATE", func(t *testing.T) {
		expected := "81"
		stack := op.NewStack()
		stack, err := op.OP1NEGATE(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ := stack.Pop()
		if instruction.Hex() != expected {
			t.Errorf("expected: %v, got: %v", expected, instruction.Hex())
		}
	})

	t.Run("OP1 to OP16", func(t *testing.T) {
		expected := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "0a", "0b", "0c", "0d", "0e", "0f", "10"}
		ops := []func(*op.Stack) (*op.Stack, error){
			op.OP1, op.OP2, op.OP3, op.OP4, op.OP5, op.OP6, op.OP7, op.OP8, op.OP9, op.OP10, op.OP11, op.OP12, op.OP13, op.OP14, op.OP15, op.OP16,
		}

		for i, opX := range ops {
			stack := op.NewStack()
			stack, err := opX(stack)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			if stack.Size() != 1 {
				t.Errorf("expected: %v, got: %v", 1, stack.Size())
			}

			instruction, _ := stack.Pop()
			if instruction.Hex() != expected[i] {
				t.Errorf("expected: %v, got: %v", expected[i], instruction.Hex())
			}
		}
	})
}

func TestNOP(t *testing.T) {
	instruction, _ := op.NewInstruction([]byte{0x01})
	stack := op.NewStack()
	stack.Push(instruction)

	stack, err := op.NOP(stack)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if stack.Size() != 1 {
		t.Errorf("expected: %v, got: %v", 1, stack.Size())
	}
}

func TestIF(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, _, err := op.IF(stack, []op.Instruction{})
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected 'stack too small', got %v", err)
		}
	})

	t.Run("If missing OP_ENDIF", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(instruction)

		_, _, err := op.IF(stack, []op.Instruction{})
		if err == nil || err.Error() != "missing OP_ENDIF" {
			t.Errorf("expected 'missing OP_ENDIF', got %v", err)
		}
	})
}

func TestVERIFY(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.VERIFY(stack)
		if err == nil || err.Error() != "transaction invalid" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Verify valid transaction", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(instruction)

		stack, err := op.VERIFY(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 0 {
			t.Errorf("expected: %v, got: %v", 0, stack.Size())
		}
	})

	t.Run("Verify invalid transaction", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x00})
		stack := op.NewStack()
		stack.Push(instruction)

		_, err := op.VERIFY(stack)
		if err == nil || err.Error() != "transaction invalid" {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestRETURN(t *testing.T) {
	stack := op.NewStack()
	stack, err := op.RETURN(stack)
	if err == nil || err.Error() != "transaction invalid" {
		t.Errorf("expected error, got nil")
	}

	if stack.Size() != 0 {
		t.Errorf("expected: %v, got: %v", 0, stack.Size())
	}
}

func TestPICK(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.PICK(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Pick element", func(t *testing.T) {
		one, _ := op.NewInstruction([]byte{0x01})
		two, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()

		stack.Push(two)
		stack.Push(one)
		stack.Push(one)

		stack, err := op.PICK(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 3 {
			t.Errorf("expected: %v, got: %v", 3, stack.Size())
		}

		top, _ := stack.Peek()
		if !top.Equals(two) {
			t.Errorf("expected: %v, got: %v", two, top)
		}
	})
}

func TestSWAP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.SWAP(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Swap elements", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.SWAP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 2 {
			t.Errorf("expected: %v, got: %v", 2, stack.Size())
		}

		element3, _ := stack.Pop()
		if !element3.Equals(element1) {
			t.Errorf("expected: %v, got: %v", element1, element3)
		}

		element4, _ := stack.Pop()
		if !element4.Equals(element2) {
			t.Errorf("expected: %v, got: %v", element2, element4)
		}
	})
}

func TestSIZE(t *testing.T) {
	t.Run("Stack too small", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.SIZE(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Size", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.SIZE(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 3 {
			t.Errorf("expected: %v, got: %v", 3, stack.Size())
		}

		lengthOfElement2, _ := stack.Pop()
		if lengthOfElement2.Hex() != "01" {
			t.Errorf("expected: %v, got: %v", "02", lengthOfElement2.Hex())
		}
	})
}

func TestOP2DUP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.OP2DUP(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Duplicate elements", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.OP2DUP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 4 {
			t.Errorf("expected: %v, got: %v", 4, stack.Size())
		}

		element3, _ := stack.PeekN(2)
		if element3.Equals(element1) {
			t.Errorf("expected: %v, got: %v", element1, element3)
		}

		element4, _ := stack.PeekN(3)
		if element4.Equals(element2) {
			t.Errorf("expected: %v, got: %v", element2, element4)
		}
	})
}

func TestOP2DROP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.OP2DROP(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Drop elements", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.OP2DROP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 0 {
			t.Errorf("expected: %v, got: %v", 0, stack.Size())
		}
	})
}

func TestIFDUP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.IFDUP(stack)
		if err == nil || err.Error() != "stack is empty" {
			t.Errorf("expected error, got %v", err)
		}
	})

	t.Run("Duplicate element", func(t *testing.T) {
		element, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(element)

		stack, err := op.IFDUP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 2 {
			t.Errorf("expected: %v, got: %v", 2, stack.Size())
		}

		element2, _ := stack.Peek()
		if !element2.Equals(element) {
			t.Errorf("expected: %v, got: %v", element, element2)
		}
	})

	t.Run("Do not duplicate element", func(t *testing.T) {
		element, _ := op.NewInstruction([]byte{0x00})
		stack := op.NewStack()
		stack.Push(element)

		stack, err := op.IFDUP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}
	})
}

func TestDEPTH(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		stack, err := op.DEPTH(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}
	})

	t.Run("Get stack depth", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.DEPTH(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 3 {
			t.Errorf("expected: %v, got: %v", 3, stack.Size())
		}

		instruction, _ := stack.Pop()
		if instruction.Hex() != "02" {
			t.Errorf("expected: %v, got: %v", "02", instruction.Hex())
		}
	})
}

func TestDROP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.DROP(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Drop instruction", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(instruction)

		stack, err := op.DROP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 0 {
			t.Errorf("expected: %v, got: %v", 0, stack.Size())
		}
	})
}

func TestDUP(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.DUP(stack)
		if err == nil || err.Error() != "stack is empty" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Duplicate instruction", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(instruction)

		stack, err := op.DUP(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 2 {
			t.Errorf("expected: %v, got: %v", 2, stack.Size())
		}
	})
}

func TestEQUAL(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.EQUAL(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Equal elements", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.EQUAL(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ := stack.Pop()
		if instruction.Hex() != "01" {
			t.Errorf("expected: %v, got: %v", "01", instruction.Hex())
		}
	})

	t.Run("Different elements", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.EQUAL(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ := stack.Pop()
		if instruction.Hex() != "00" {
			t.Errorf("expected: %v, got: %v", "00", instruction.Hex())
		}
	})
}

func TestNOT(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.NOT(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Not instruction on 0x42 value", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x42})
		stack := op.NewStack()
		stack.Push(instruction)

		stack, err := op.NOT(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ = stack.Pop()
		if instruction.Hex() != "00" {
			t.Errorf("expected: %v, got: %v", "00", instruction.Hex())
		}
	})

	t.Run("Not instruction on 0x00 value", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x00})
		stack := op.NewStack()
		stack.Push(instruction)

		stack, err := op.NOT(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ = stack.Pop()
		if instruction.Hex() != "01" {
			t.Errorf("expected: %v, got: %v", "01", instruction.Hex())
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
		element1, _ := op.NewInstruction([]byte{0x01})
		element2, _ := op.NewInstruction([]byte{0x02})
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

		instruction, _ := stack.Pop()
		if instruction.Hex() != "03" {
			t.Errorf("expected: %v, got: %v", "03", instruction.Hex())
		}
	})
}

func TestMUL(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.MUL(stack)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Multiply elements", func(t *testing.T) {
		element1, _ := op.NewInstruction([]byte{0x02})
		element2, _ := op.NewInstruction([]byte{0x03})
		stack := op.NewStack()
		stack.Push(element1)
		stack.Push(element2)

		stack, err := op.MUL(stack)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if stack.Size() != 1 {
			t.Errorf("expected: %v, got: %v", 1, stack.Size())
		}

		instruction, _ := stack.Pop()
		if instruction.Hex() != "06" {
			t.Errorf("expected: %v, got: %v", "06", instruction.Hex())
		}
	})
}

func TestSHA1(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.SHA1(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Hash instruction", func(t *testing.T) {
		expected := "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
		instruction, _ := op.NewInstruction([]byte("hello world"))
		stack := op.NewStack()
		stack.Push(instruction)

		stack, err := op.SHA1(stack)
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

func TestHASH160(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.HASH160(stack)
		if err == nil || err.Error() != "invalid stack" {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("Hash instruction", func(t *testing.T) {
		expected := "d7d5ee7824ff93f94c3055af9382c86c68b5ca92"
		instruction, _ := op.NewInstruction([]byte("hello world"))
		stack := op.NewStack()
		stack.Push(instruction)

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

	t.Run("Hash instruction", func(t *testing.T) {
		instruction, _ := op.NewInstruction([]byte{0x01})
		stack := op.NewStack()
		stack.Push(instruction)

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

		sec, _ := op.NewInstruction(secBytes)
		sig, _ := op.NewInstruction(sigBytes)

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

		instruction, _ := stack.Pop()
		if instruction.Hex() != "01" {
			t.Errorf("expected: %v, got: %v", "01", instruction.Hex())
		}
	})
}

func TestCHECKMULTISIG(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		stack := op.NewStack()
		_, err := op.CHECKMULTISIG(stack, nil)
		if err == nil || err.Error() != "stack too small" {
			t.Errorf("expected error, got nil")
		}
	})
}
