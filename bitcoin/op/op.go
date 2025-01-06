package op

import (
	"fmt"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
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

// The number 1 is pushed onto the stack.
func OP1(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x01})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// 16 is pushed onto the stack.
func OP16(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x10})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// Duplicates the top stack item.
func DUP(stack *Stack) (*Stack, error) {
	duplicateElement, err := stack.Peek()
	if err != nil {
		return nil, err
	}

	stack.Push(duplicateElement)

	return stack, nil
}

// a is added to b.
func ADD(stack *Stack) (*Stack, error) {
	if stack.Size() < 2 {
		return nil, fmt.Errorf("stack too small")
	}

	b, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	a, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	bInt := new(big.Int).SetBytes(b.element)
	aInt := new(big.Int).SetBytes(a.element)

	sum := new(big.Int).Add(aInt, bInt)
	sumElement, err := NewElement(sum.Bytes())
	if err != nil {
		return nil, err
	}

	stack.Push(sumElement)

	return stack, nil
}

// The input is hashed twice: first with SHA-256 and then with RIPEMD-160.
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

// The input is hashed two times with SHA-256.
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

// The entire transaction's outputs, inputs, and script (from the
// most recently-executed OP_CODESEPARATOR to the end) are hashed.
// The signature used by OP_CHECKSIG must be a valid signature for
// this hash and public key. If it is, 1 is returned, 0 otherwise.
func CHECKSIG(stack *Stack, z *big.Int) (*Stack, error) {
	if stack.Size() < 2 {
		return nil, fmt.Errorf("stack too small")
	}

	secPubKey, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	derSignature, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	point, err := ecc.Parse(secPubKey.element)
	if err != nil {
		return nil, err
	}

	signature, err := ecc.ParseDER(derSignature.element[:len(derSignature.element)-1])
	if err != nil {
		return nil, err
	}

	valid, err := point.Verify(z, signature)
	if err != nil {
		return nil, err
	}

	var data []byte
	if valid {
		data = []byte{0x01}
	} else {
		data = []byte{0x00}
	}

	newElement, err := NewElement(data)
	if err != nil {
		return nil, err
	}

	stack.Push(newElement)

	return stack, nil
}

var OP_CODE_FUNCTIONS = map[int]func(*Stack) (*Stack, error){
	0:   OP0,
	81:  OP1,
	96:  OP16,
	118: DUP,
	147: ADD,
	169: HASH160,
	170: HASH256,
	// 172: CHECKSIG,
}
