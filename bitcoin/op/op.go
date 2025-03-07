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
	instruction, err := NewInstruction([]byte{})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number -1 is pushed onto the stack.
func OP1NEGATE(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x81})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 1 is pushed onto the stack.
func OP1(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x01})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 2 is pushed onto the stack.
func OP2(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x02})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 3 is pushed onto the stack.
func OP3(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x03})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 4 is pushed onto the stack.
func OP4(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x04})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 5 is pushed onto the stack.
func OP5(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x05})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 6 is pushed onto the stack.
func OP6(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x06})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 7 is pushed onto the stack.
func OP7(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x07})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 8 is pushed onto the stack.
func OP8(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x08})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 9 is pushed onto the stack.
func OP9(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x09})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 10 is pushed onto the stack.
func OP10(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x0A})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 11 is pushed onto the stack.
func OP11(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x0B})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 12 is pushed onto the stack.
func OP12(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x0C})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 13 is pushed onto the stack.
func OP13(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x0D})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 14 is pushed onto the stack.
func OP14(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x0E})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// The number 15 is pushed onto the stack.
func OP15(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x0F})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// 16 is pushed onto the stack.
func OP16(stack *Stack) (*Stack, error) {
	instruction, err := NewInstruction([]byte{0x10})
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// Does nothing.
func NOP(stack *Stack) (*Stack, error) {
	return stack, nil
}

// If the top stack value is not False, the statements are executed. The top stack value is removed.
func IF(stack *Stack, instructions []Instruction) (*Stack, *[]Instruction, error) {
	if stack.Size() < 1 {
		return nil, nil, fmt.Errorf("stack too small")
	}

	founded := false
	numberOfEndIfsNeeded := 1
	trueItems := make([]Instruction, 0)
	falseItems := make([]Instruction, 0)
	current := make([]Instruction, 0)

	for _, item := range instructions {
		if bytes.Equal(item.instruction, []byte{0x63}) || bytes.Equal(item.instruction, []byte{0x64}) {
			numberOfEndIfsNeeded++
			current = append(current, item)
		} else if numberOfEndIfsNeeded == 1 || numberOfEndIfsNeeded == 103 {
			current = falseItems
		} else if bytes.Equal(item.instruction, []byte{0x68}) {
			if numberOfEndIfsNeeded == 1 {
				founded = true
				break
			} else {
				numberOfEndIfsNeeded--
				current = append(current, item)
			}
		} else {
			current = append(current, item)
		}
	}

	if !founded {
		return nil, nil, fmt.Errorf("missing OP_ENDIF")
	}

	instruction, err := stack.Pop()
	if err != nil {
		return nil, nil, err
	}

	if bytes.Equal(instruction.instruction, []byte{0x00}) {
		instructions = append(falseItems, instructions...)
	} else {
		instructions = append(trueItems, instructions...)
	}

	return stack, &instructions, nil
}

// Marks transaction as invalid if top stack value is not true. The top stack value is removed.
func VERIFY(stack *Stack) (*Stack, error) {
	if stack.Size() < 1 {
		return nil, fmt.Errorf("transaction invalid") // Should stack be included in the return?
	}

	instruction, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	if bytes.Equal(instruction.instruction, []byte{0x00}) {
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

// Removes the top two stack items.
func OP2DROP(stack *Stack) (*Stack, error) {
	if stack.Size() < 2 {
		return nil, fmt.Errorf("stack too small")
	}

	_, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	_, err = stack.Pop()
	if err != nil {
		return nil, err
	}

	return stack, nil
}

// If the top stack value is not 0, duplicate it.
func IFDUP(stack *Stack) (*Stack, error) {
	instruction, err := stack.Peek()
	if err != nil {
		return nil, err
	}

	if instruction.IsZero() {
		return stack, nil
	}

	stack.Push(instruction)

	return stack, nil
}

// Puts the number of stack items onto the stack.
func DEPTH(stack *Stack) (*Stack, error) {
	size := []byte{byte(stack.Size())}
	instruction, err := NewInstruction(size)
	if err != nil {
		return nil, err
	}

	stack.Push(instruction)

	return stack, nil
}

// Removes the top stack item.
func DROP(stack *Stack) (*Stack, error) {
	_, err := stack.Pop()
	if err != nil {
		return nil, err
	}

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

// The item n back in the stack is copied to the top.
func PICK(stack *Stack) (*Stack, error) {
	if stack.Size() < 1 {
		return nil, fmt.Errorf("stack too small")
	}

	data, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	index := new(big.Int).Add(new(big.Int).SetInt64(1), new(big.Int).SetBytes(data.instruction))
	if index.Cmp(big.NewInt(int64(stack.Size()))) > 0 {
		return nil, fmt.Errorf("stack too small")
	}

	element, err := stack.PeekN(int(index.Int64()) - 1)
	if err != nil {
		return nil, err
	}

	stack.Push(element)

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

// Pushes the string length of the top element of the stack (without popping it).
func SIZE(stack *Stack) (*Stack, error) {
	if stack.Size() < 1 {
		return nil, fmt.Errorf("stack too small")
	}

	instruction, err := stack.Peek()
	if err != nil {
		return nil, err
	}

	size := len(instruction.instruction)
	data, err := NewInstruction([]byte{byte(size)})
	if err != nil {
		return nil, err
	}

	stack.Push(data)

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
	if bytes.Equal(element1.instruction, element2.instruction) {
		data = []byte{0x01}
	} else {
		data = []byte{0x00}
	}

	NewInstruction, err := NewInstruction(data)
	if err != nil {
		return nil, err
	}

	stack.Push(NewInstruction)

	return stack, nil
}

// If the input is 0 or 1, it is flipped. Otherwise the output will be 0.
func NOT(stack *Stack) (*Stack, error) {
	instruction, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	value := new(big.Int).SetBytes(instruction.instruction)

	var data []byte
	if value.Cmp(big.NewInt(0)) == 0 {
		data = []byte{0x01}
	} else {
		data = []byte{0x00}
	}

	NewInstruction, err := NewInstruction(data)
	if err != nil {
		return nil, err
	}

	stack.Push(NewInstruction)

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

	bInt := new(big.Int).SetBytes(b.instruction)
	aInt := new(big.Int).SetBytes(a.instruction)

	sum := new(big.Int).Add(aInt, bInt)
	sumElement, err := NewInstruction(sum.Bytes())
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

	bInt := new(big.Int).SetBytes(b.instruction)
	aInt := new(big.Int).SetBytes(a.instruction)

	product := new(big.Int).Mul(aInt, bInt)
	productElement, err := NewInstruction(product.Bytes())
	if err != nil {
		return nil, err
	}

	stack.Push(productElement)

	return stack, nil
}

// Returns the smaller of a and b.
func MIN(stack *Stack) (*Stack, error) {
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
	if element1.Int64() < element2.Int64() {
		stack.Push(element1)
	} else {
		stack.Push(element2)
	}

	return stack, nil
}

// Returns the larger of a and b.
func MAX(stack *Stack) (*Stack, error) {
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
	if element1.Int64() > element2.Int64() {
		stack.Push(element1)
	} else {
		stack.Push(element2)
	}

	return stack, nil
}

// Returns 1 if x is within the specified range (left-inclusive), 0 otherwise.
func WITHIN(stack *Stack) (*Stack, error) {
	if stack.Size() < 3 {
		return nil, fmt.Errorf("stack too small")
	}

	max, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	min, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	x, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	if x.Int64() >= min.Int64() && x.Int64() < max.Int64() {
		one, err := NewInstruction([]byte{0x01})
		if err != nil {
			return nil, err
		}
		stack.Push(one)
	} else {
		zero, err := NewInstruction([]byte{0x00})
		if err != nil {
			return nil, err
		}
		stack.Push(zero)
	}

	return stack, nil
}

// The input is hashed using SHA-1.
func SHA1(stack *Stack) (*Stack, error) {
	instruction, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	hashed := hash.HashSHA1(instruction.instruction)
	hashedElement, err := NewInstruction(hashed)
	if err != nil {
		return nil, err
	}
	stack.Push(hashedElement)

	return stack, nil
}

// The input is hashed twice: first with SHA-256 and then with RIPEMD-160.
func HASH160(stack *Stack) (*Stack, error) {
	instruction, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	hashed := hash.Hash160(instruction.instruction)
	hashedElement, err := NewInstruction(hashed)
	if err != nil {
		return nil, err
	}
	stack.Push(hashedElement)

	return stack, nil
}

// The input is hashed two times with SHA-256.
func HASH256(stack *Stack) (*Stack, error) {
	instruction, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	hashed := hash.Hash256(instruction.instruction)
	hashedElement, err := NewInstruction(hashed)
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

	point, err := ecc.Parse(secPubKey.instruction)
	if err != nil {
		return nil, err
	}

	signature, err := ecc.ParseDER(derSignature.instruction[:len(derSignature.instruction)-1])
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

	NewInstruction, err := NewInstruction(data)
	if err != nil {
		return nil, err
	}

	stack.Push(NewInstruction)

	return stack, nil
}

// Compares the first signature against each public key until it
// finds an ECDSA match. Starting with the subsequent public key,
// it compares the second signature against each remaining public
// key until it finds an ECDSA match. The process is repeated
// until all signatures have been checked or not enough public
// keys remain to produce a successful result. All signatures
// need to match a public key. Because public keys are not checked
// again if they fail any signature comparison, signatures must be
// placed in the scriptSig using the same order as their
// corresponding public keys were placed in the scriptPubKey or
// redeemScript. If all signatures are valid, 1 is returned, 0
// otherwise. Due to a bug, one extra unused value is removed from
// the stack.
func CHECKMULTISIG(stack *Stack, z *big.Int) (*Stack, error) {
	if stack.Size() < 1 {
		return nil, fmt.Errorf("stack too small")
	}

	n, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	if int64(stack.Size()) < (n.Int64() + 1) {
		return nil, fmt.Errorf("stack too small")
	}

	secPubKey := make([]Instruction, n.Int64())

	for i := 0; i < int(n.Int64()); i++ {
		instruction, err := stack.Pop()
		if err != nil {
			return nil, err
		}

		secPubKey[i] = *instruction
	}

	m, err := stack.Pop()
	if err != nil {
		return nil, err
	}

	derSignatures := make([]Instruction, m.Int64())

	for i := 0; i < int(m.Int64()); i++ {
		instruction, err := stack.Pop()
		if err != nil {
			return nil, err
		}

		derSignatures[i] = *instruction
	}

	_, err = stack.Pop()
	if err != nil {
		return nil, err
	}

	points := make([]*ecc.S256Point, len(secPubKey))
	for index, sec := range secPubKey {
		point, err := ecc.Parse(sec.Bytes())
		if err != nil {
			return nil, err
		}

		points[index] = point
	}

	sigs := make([]*ecc.Signature, len(derSignatures))
	for index, der := range derSignatures {
		signature, err := ecc.ParseDER(der.Bytes())
		if err != nil {
			return nil, err
		}

		sigs[index] = signature
	}

	for _, sig := range sigs {
		if len(points) == 0 {
			return nil, fmt.Errorf("signatures do not match public keys")
		}

		for len(points) > 0 {
			point := points[0]
			points = points[1:]

			ok, err := point.Verify(z, sig)
			if err != nil {
				return nil, err
			}

			if ok {
				break
			}
		}
	}

	data, err := NewInstruction([]byte{0x01})
	if err != nil {
		return nil, err
	}
	stack.Push(data)

	return stack, nil
}

var OP_CODE_FUNCTIONS = map[int]func(*Stack) (*Stack, error){
	0:  OP0,
	79: OP1NEGATE,
	81: OP1,
	82: OP2,
	83: OP3,
	84: OP4,
	85: OP5,
	86: OP6,
	87: OP7,
	88: OP8,
	89: OP9,
	90: OP10,
	91: OP11,
	92: OP12,
	93: OP13,
	94: OP14,
	95: OP15,
	96: OP16,
	97: NOP,
	// 99:  IF,
	105: VERIFY,
	106: RETURN,
	109: OP2DROP,
	110: OP2DUP,
	115: IFDUP,
	116: DEPTH,
	117: DROP,
	118: DUP,
	121: PICK,
	123: SWAP,
	130: SIZE,
	135: EQUAL,
	145: NOT,
	147: ADD,
	149: MUL,
	163: MIN,
	164: MAX,
	165: WITHIN,
	167: SHA1,
	169: HASH160,
	170: HASH256,
	// 172: CHECKSIG,
	// 174: CHECKMULTISIG,
}

var OP_CODE = struct {
	EQUAL   Instruction
	HASH160 Instruction
}{
	EQUAL:   Instruction{instruction: []byte{0x87}},
	HASH160: Instruction{instruction: []byte{0xa9}},
}
