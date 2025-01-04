package op

import (
	"github.com/stefanalfbo/programmingbitcoin/ecc"
)

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
	118: DUP,
	169: HASH160,
	170: HASH256,
}
