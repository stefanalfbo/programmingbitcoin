package bitcoin

import (
	"fmt"
	"io"
	"math/big"
	"slices"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin/op"
	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type Instruction struct {
	instruction []byte
}

func NewInstruction(instruction []byte) *Instruction {
	return &Instruction{instruction}
}

func (i *Instruction) Hex() string {
	return fmt.Sprintf("%x", i.instruction)
}

type Script struct {
	instructions []Instruction
}

func NewScript(instructions []Instruction) *Script {
	return &Script{instructions}
}

func (script *Script) String() string {
	return "notImplementedYet"
}

func (script *Script) Instructions() []Instruction {
	return script.instructions
}

func ParseScript(data io.Reader) (*Script, error) {
	length, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	instructions := make([]Instruction, 0)
	scriptLength := int(length.Int64())

	count := 0
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
			instructions = append(instructions, Instruction{tmpData})
			count += int(currentByte)
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

			instructions = append(instructions, Instruction{tmpData})

			count += int(dataLength) + 1
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

			instructions = append(instructions, Instruction{tmpData})

			count += int(dataLength) + 2
		} else {
			opCode := []byte{currentByte}
			instructions = append(instructions, Instruction{opCode})
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

		if isOpCode(instruction.instruction) {
			scriptAsBytes = append(scriptAsBytes, instruction.instruction...)
		} else {
			length := len(instruction.instruction)

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
			scriptAsBytes = append(scriptAsBytes, instruction.instruction...)
		}
	}

	return scriptAsBytes, nil
}

func isOpCode(opCode []byte) bool {
	if len(opCode) != 1 {
		return false
	}
	return opCode[0] >= 0x01 && opCode[0] <= 0x4b
}

func (script *Script) Serialize() ([]byte, error) {
	scriptAsBytes, err := script.RawSerialize()
	if err != nil {
		return nil, err
	}

	scriptLength := len(scriptAsBytes)
	length, err := varint.Encode(big.NewInt(int64(scriptLength)))
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

	for _, instruction := range script.instructions {
		if isOpCode(instruction.instruction) {
			opCode := int(instruction.instruction[0])
			operation, exists := op.OP_CODE_FUNCTIONS[opCode]
			if exists {
				s, err := operation(stack)
				if err != nil {
					return false, err
				}
				stack = s
			} else {
				if opCode == 99 {
					// s, err := op.IF(stack, script)
					// if err != nil {
					// 	return false, err
					// }
					// stack = s
				}
				if slices.Contains([]int{99, 100}, opCode) {
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
			element, err := op.NewElement(instruction.instruction)
			if err != nil {
				return false, err
			}
			stack.Push(element)
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
