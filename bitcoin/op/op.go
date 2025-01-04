package op

import (
	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

// An empty array of bytes is pushed onto the stack. (This is not a no-op: an item is added to the stack.)
func OP0(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

func DUP(stack *Stack) (*Stack, error) {
	duplicateElement, err := stack.Peek()
	if err != nil {
		return nil, err
	}

	stack.Push(duplicateElement)

	return stack, nil
}

func HASH160(stack *Stack) (*Stack, error) {
	element, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	hashed := ecc.Hash160(element.element)
	hashedElement, err := NewElement(hashed)
	if err != nil {
		return nil, err
	}
	stack.Push(hashedElement)

	return stack, nil
}

func HASH256(stack *Stack) (*Stack, error) {
	element, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	hashed := ecc.Hash256(string(element.element))
	hashedElement, err := NewElement(hashed.Bytes())
	if err != nil {
		return nil, err
	}
	stack.Push(hashedElement)

	return stack, nil
}

var OP_CODE_FUNCTIONS = map[int]func(*Stack) (*Stack, error){
	0:   OP0,
	118: DUP,
	169: HASH160,
	170: HASH256,
}
