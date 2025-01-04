package op

import "fmt"

type Element struct {
	element []byte
}

func NewElement(data []byte) (*Element, error) {
	elementSize := len(data)

	if elementSize > 520 {
		return nil, fmt.Errorf("element too large")
	}

	return &Element{data}, nil
}

func (e *Element) Equals(other *Element) bool {
	if len(e.element) != len(other.element) {
		return false
	}

	for i := range e.element {
		if e.element[i] != other.element[i] {
			return false
		}
	}

	return true
}

func (e *Element) Hex() string {
	return fmt.Sprintf("%x", e.element)
}

type Stack struct {
	stack []Element
}

func NewStack() *Stack {
	return &Stack{make([]Element, 0)}
}

func (s *Stack) Push(element *Element) {
	s.stack = append(s.stack, *element)
}

func (s *Stack) Size() int {
	return len(s.stack)
}

func (s *Stack) Pop() (*Element, error) {
	if s.Size() < 1 {
		return nil, fmt.Errorf("invalid stack")
	}

	element := s.stack[s.Size()-1]
	s.stack = s.stack[:s.Size()-1]

	return &element, nil
}

func (s *Stack) Peek() (*Element, error) {
	if s.Size() < 1 {
		return nil, fmt.Errorf("invalid stack")
	}

	return &s.stack[s.Size()-1], nil
}
