package bitcoin

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type Script struct {
	instructions []op.Instruction
}

func NewScript(instructions []op.Instruction) *Script {
	return &Script{instructions}
}

func (script *Script) String() string {
	return "notImplementedYet"
}

func (script *Script) Instructions() []op.Instruction {
	return script.instructions
}

// Returns whether this follows the
// OP_HASH160 <20 byte hash> OP_EQUAL pattern.
func (script *Script) IsP2SHScriptPubKey() bool {
	return len(script.instructions) == 3 &&
		script.instructions[0].Equals(&op.OP_CODE.HASH160) &&
		script.instructions[1].Length() == 20 &&
		script.instructions[2].Equals(&op.OP_CODE.EQUAL)
}

func ParseScript(data io.Reader) (*Script, error) {
	length, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	instructions := make([]op.Instruction, 0)
	scriptLength := length

	var count uint64 = 0
	for count < scriptLength {
		current := make([]byte, 1)
		_, err := data.Read(current)
		if err != nil {
			return nil, err
		}

		count += 1

		currentByte := current[0]
		if currentByte >= 1 && currentByte <= 75 {
			tmpData := make([]byte, currentByte)
			_, err := data.Read(tmpData)
			if err != nil {
				return nil, err
			}
			instruction, err := op.NewInstruction(tmpData)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instruction)
			count += uint64(currentByte)
		} else if currentByte == 76 {
			lengthContainer := make([]byte, 1)
			_, err := data.Read(lengthContainer)
			if err != nil {
				return nil, err
			}
			dataLength := endian.LittleEndianToInt32(lengthContainer)

			tmpData := make([]byte, dataLength)
			_, err = data.Read(tmpData)
			if err != nil {
				return nil, err
			}
			instruction, err := op.NewInstruction(tmpData)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instruction)

			count += uint64(dataLength) + 1
		} else if currentByte == 77 {
			lengthContainer := make([]byte, 2)
			_, err := data.Read(lengthContainer)
			if err != nil {
				return nil, err
			}
			dataLength := endian.LittleEndianToInt32(lengthContainer)

			tmpData := make([]byte, dataLength)
			_, err = data.Read(tmpData)
			if err != nil {
				return nil, err
			}
			instruction, err := op.NewInstruction(tmpData)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instruction)

			count += uint64(dataLength) + 2
		} else {
			instruction, err := op.NewInstruction([]byte{currentByte})
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instruction)
		}
	}

	if count != scriptLength {
		return nil, fmt.Errorf("parsing script failed")
	}

	return NewScript(instructions), nil
}

func (script *Script) RawSerialize() ([]byte, error) {
	scriptAsBytes := make([]byte, 0)

	for _, instruction := range script.instructions {

		if instruction.IsOpCode() {
			scriptAsBytes = append(scriptAsBytes, instruction.Bytes()...)
		} else {
			length := instruction.Length()

			if length < 76 {
				value := endian.BigIntToLittleEndian(big.NewInt(int64(length)), 1)
				scriptAsBytes = append(scriptAsBytes, value...)
			} else if length < 256 {
				value := endian.BigIntToLittleEndian(big.NewInt(int64(76)), 2)
				scriptAsBytes = append(scriptAsBytes, value...)

				value = endian.BigIntToLittleEndian(big.NewInt(int64(length)), 1)
				scriptAsBytes = append(scriptAsBytes, value...)
			} else if length <= 520 {
				value := endian.BigIntToLittleEndian(big.NewInt(int64(77)), 2)
				scriptAsBytes = append(scriptAsBytes, value...)

				value = endian.BigIntToLittleEndian(big.NewInt(int64(length)), 2)
				scriptAsBytes = append(scriptAsBytes, value...)
			} else {
				return nil, fmt.Errorf("too long of an instruction")
			}
			scriptAsBytes = append(scriptAsBytes, instruction.Bytes()...)
		}
	}

	return scriptAsBytes, nil
}

func (script *Script) Serialize() ([]byte, error) {
	scriptAsBytes, err := script.RawSerialize()
	if err != nil {
		return nil, err
	}

	scriptLength := uint64(len(scriptAsBytes))
	length, err := varint.Encode(scriptLength)
	if err != nil {
		return nil, err
	}

	return append(length, scriptAsBytes...), nil
}

func (script *Script) Add(other *Script) *Script {
	instructions := append(script.instructions, other.instructions...)
	return NewScript(instructions)
}

func (script *Script) Evaluate(z []byte) (bool, error) {

	stack := op.NewStack()
	// altStack := NewStack()

	for len(script.instructions) > 0 {
		instruction := script.instructions[0]
		script.instructions = script.instructions[1:]

		if instruction.IsOpCode() {
			opCode := int(instruction.Bytes()[0])
			operation, exists := op.OP_CODE_FUNCTIONS[opCode]
			if exists {
				s, err := operation(stack)
				if err != nil {
					return false, err
				}
				stack = s
			} else {
				if opCode == 99 {
					s, instructions, err := op.IF(stack, script.instructions)
					if err != nil {
						return false, err
					}
					stack = s
					script.instructions = *instructions
				} else if opCode == 100 {
					// OP_IF, OP_NOTIF

				} else if slices.Contains([]int{107, 108}, opCode) {
					// OP_TOALTSTACK, OP_FROMALTSTACK
				} else if slices.Contains([]int{172, 173, 174, 175}, opCode) {
					// OP_CHECKSIG, OP_CHECKSIGVERIFY, OP_CHECKMULTISIG, OP_CHECKMULTISIGVERIFY
					s, err := op.CHECKSIG(stack, new(big.Int).SetBytes(z))
					if err != nil {
						return false, err
					}
					stack = s
				}
			}
		} else {
			element, err := op.NewInstruction(instruction.Bytes())
			if err != nil {
				return false, err
			}
			stack.Push(element)
			if len(script.instructions) == 3 &&
				script.instructions[0].Equals(&op.OP_CODE.HASH160) &&
				script.instructions[1].Length() == 20 &&
				script.instructions[2].Equals(&op.OP_CODE.EQUAL) {

				hash160 := script.instructions[1]

				script.instructions = script.instructions[3:]
				stack, err := op.HASH160(stack)
				if err != nil {
					return false, err
				}
				stack.Push(&hash160)

				stack, err = op.EQUAL(stack)
				if err != nil {
					return false, err
				}

				_, err = op.VERIFY(stack)
				if err != nil {
					return false, err
				}

				length, err := varint.Encode(uint64(instruction.Length()))
				if err != nil {
					return false, err
				}

				redeemScript, err := ParseScript(bytes.NewReader(append(length, instruction.Bytes()...)))
				if err != nil {
					return false, err
				}

				script.instructions = append(redeemScript.instructions, script.instructions...)
			}
		}
	}

	if stack.Size() == 0 {
		return false, fmt.Errorf("stack size is not 1")
	}

	element, err := stack.Pop()
	if err != nil {
		return false, err
	}

	if element == nil {
		return false, fmt.Errorf("element is nil")
	}
	return true, nil
}

func ToP2PKHScript(h160 []byte) (*Script, error) {

	dup, err := op.NewInstruction([]byte{0x76})
	if err != nil {
		return nil, err
	}

	hash160, err := op.NewInstruction([]byte{0xa9})
	if err != nil {
		return nil, err
	}

	data, err := op.NewInstruction(h160)
	if err != nil {
		return nil, err
	}

	equalVerify, err := op.NewInstruction([]byte{0x88})
	if err != nil {
		return nil, err
	}

	checkSignature, err := op.NewInstruction([]byte{0xac})
	if err != nil {
		return nil, err
	}

	p2pkh := NewScript([]op.Instruction{
		*dup,
		*hash160,
		*data,
		*equalVerify,
		*checkSignature,
	})

	return p2pkh, nil
}
