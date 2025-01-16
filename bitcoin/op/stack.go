package op

import (
	"fmt"
	"math/big"
)

// Instruction represents a single instruction or a value in a script.
type Instruction struct {
	instruction []byte
}

func NewInstruction(instruction []byte) (*Instruction, error) {
	instructionSize := len(instruction)

	if instructionSize > 520 {
		return nil, fmt.Errorf("instruction too large")
	}

	return &Instruction{instruction}, nil
}

func (i *Instruction) Hex() string {
	return fmt.Sprintf("%x", i.instruction)
}

func (e *Instruction) Equals(other *Instruction) bool {
	if len(e.instruction) != len(other.instruction) {
		return false
	}

	for i := range e.instruction {
		if e.instruction[i] != other.instruction[i] {
			return false
		}
	}

	return true
}

func (i *Instruction) IsZero() bool {
	for _, b := range i.instruction {
		if b != 0 {
			return false
		}
	}

	return true
}

func (i *Instruction) IsOpCode() bool {
	if len(i.instruction) != 1 {
		return false
	}

	return i.instruction[0] >= 0x01 && i.instruction[0] <= 0x4b
}

func (i *Instruction) Length() int {
	return len(i.instruction)
}

func (i *Instruction) Bytes() []byte {
	return i.instruction
}

func (i *Instruction) Int64() int64 {
	n := big.NewInt(0).SetBytes(i.instruction)
	return n.Int64()
}

type Stack struct {
	stack []Instruction
}

func NewStack() *Stack {
	return &Stack{make([]Instruction, 0)}
}

func (s *Stack) Push(element *Instruction) {
	s.stack = append(s.stack, *element)
}

func (s *Stack) Size() int {
	return len(s.stack)
}

func (s *Stack) Pop() (*Instruction, error) {
	if s.Size() < 1 {
		return nil, fmt.Errorf("invalid stack")
	}

	element := s.stack[s.Size()-1]
	s.stack = s.stack[:s.Size()-1]

	return &element, nil
}

func (s *Stack) Peek() (*Instruction, error) {
	if s.Size() < 1 {
		return nil, fmt.Errorf("invalid stack")
	}

	return &s.stack[s.Size()-1], nil
}

func (s *Stack) PeekN(n int) (*Instruction, error) {
	if s.Size() < n+1 {
		return nil, fmt.Errorf("invalid stack")
	}

	return &s.stack[s.Size()-n-1], nil
}
