package bitcoin

import (
	"fmt"
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

type Instruction struct {
	instruction []byte
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
