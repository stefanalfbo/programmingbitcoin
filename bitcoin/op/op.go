package op

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/stefanalfbo/programmingbitcoin/crypto/ecc"
	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
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

// The number -1 is pushed onto the stack.
func OP1NEGATE(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x81})
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

// The number 2 is pushed onto the stack.
func OP2(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x02})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 3 is pushed onto the stack.
func OP3(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x03})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 4 is pushed onto the stack.
func OP4(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x04})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 5 is pushed onto the stack.
func OP5(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x05})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 6 is pushed onto the stack.
func OP6(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x06})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 7 is pushed onto the stack.
func OP7(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x07})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 8 is pushed onto the stack.
func OP8(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x08})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 9 is pushed onto the stack.
func OP9(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x09})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 10 is pushed onto the stack.
func OP10(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x0A})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 11 is pushed onto the stack.
func OP11(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x0B})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 12 is pushed onto the stack.
func OP12(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x0C})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 13 is pushed onto the stack.
func OP13(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x0D})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 14 is pushed onto the stack.
func OP14(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x0E})
	if err != nil {
		return nil, err
	}

	stack.Push(element)

	return stack, nil
}

// The number 15 is pushed onto the stack.
func OP15(stack *Stack) (*Stack, error) {
	element, err := NewElement([]byte{0x0F})
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

// Does nothing.
func NOP(stack *Stack) (*Stack, error) {
	return stack, nil
}

// Marks transaction as invalid if top stack value is not true. The top stack value is removed.
func VERIFY(stack *Stack) (*Stack, error) {
	if stack.Size() < 1 {
		return nil, fmt.Errorf("transaction invalid") // Should stack be included in the return?
	}

	element, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	if bytes.Equal(element.element, []byte{0x00}) {
		return nil, fmt.Errorf("transaction invalid") // Should stack be included in the return?
	}

	return stack, nil
}

// Marks transaction as invalid. Since bitcoin 0.9, a standard way of attaching extra data to
// transactions is to add a zero-value output with a scriptPubKey consisting of OP_RETURN followed
// by data. Such outputs are provably unspendable and specially discarded from storage in the UTXO
// set, reducing their cost to the network. Since 0.12, standard relay rules allow a single output
// with OP_RETURN, that contains any sequence of push statements (or OP_RESERVED[1]) after the
// OP_RETURN provided the total scriptPubKey length is at most 83 bytes.
func RETURN(stack *Stack) (*Stack, error) {
	return stack, fmt.Errorf("transaction invalid")
}

// Duplicates the top two stack items.
func OP2DUP(stack *Stack) (*Stack, error) {
	if stack.Size() < 2 {
		return nil, fmt.Errorf("stack too small")
	}

	element1, err := stack.Peek()
	if err != nil {
		return nil, err
	}

	element2, err := stack.PeekN(1)
	if err != nil {
		return nil, err
	}

	stack.Push(element2)
	stack.Push(element1)

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

// The top two items on the stack are swapped.
func SWAP(stack *Stack) (*Stack, error) {
	if stack.Size() < 2 {
		return nil, fmt.Errorf("stack too small")
	}

	element1, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	element2, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	stack.Push(element1)
	stack.Push(element2)

	return stack, nil
}

// Returns 1 if the inputs are exactly equal, 0 otherwise.
func EQUAL(stack *Stack) (*Stack, error) {
	if stack.Size() < 2 {
		return nil, fmt.Errorf("stack too small")
	}

	element1, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	element2, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	var data []byte
	if bytes.Equal(element1.element, element2.element) {
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

// a is multiplied by b. disabled.
func MUL(stack *Stack) (*Stack, error) {
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

	product := new(big.Int).Mul(aInt, bInt)
	productElement, err := NewElement(product.Bytes())
	if err != nil {
		return nil, err
	}

	stack.Push(productElement)

	return stack, nil
}

// The input is hashed twice: first with SHA-256 and then with RIPEMD-160.
func HASH160(stack *Stack) (*Stack, error) {
	element, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	hashed := hash.Hash160(element.element)
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

	hashed := hash.Hash256(element.element)
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
	79:  OP1NEGATE,
	81:  OP1,
	82:  OP2,
	83:  OP3,
	84:  OP4,
	85:  OP5,
	86:  OP6,
	87:  OP7,
	88:  OP8,
	89:  OP9,
	90:  OP10,
	91:  OP11,
	92:  OP12,
	93:  OP13,
	94:  OP14,
	95:  OP15,
	96:  OP16,
	97:  NOP,
	105: VERIFY,
	106: RETURN,
	110: OP2DUP,
	118: DUP,
	123: SWAP,
	135: EQUAL,
	147: ADD,
	149: MUL,
	169: HASH160,
	170: HASH256,
	// 172: CHECKSIG,
}
