package varint

import (
	"encoding/binary"
	"io"
)

// Decode reads a variable uint64 from a stream.
func Decode(stream io.Reader) (uint64, error) {
	i := make([]byte, 1)

	_, err := stream.Read(i)
	if err != nil {
		return 0, err
	}

	if i[0] < 0xfd {
		return uint64(i[0]), nil
	}

	if i[0] == 0xfd {
		n := make([]byte, 2)
		_, err := stream.Read(n)
		if err != nil {
			return 0, err
		}
		n = append(n, []byte{0, 0, 0, 0, 0, 0}...)
		return binary.LittleEndian.Uint64(n), nil
	}

	if i[0] == 0xfe {
		n := make([]byte, 4)
		_, err := stream.Read(n)
		if err != nil {
			return 0, err
		}
		n = append(n, []byte{0, 0, 0, 0}...)
		return binary.LittleEndian.Uint64(n), nil
	}

	if i[0] == 0xff {
		n := make([]byte, 8)
		_, err := stream.Read(n)
		if err != nil {
			return 0, err
		}
		return binary.LittleEndian.Uint64(n), nil
	}

	return binary.LittleEndian.Uint64(i), nil
}

// Encodes a uint64 to a varint
func Encode(n uint64) ([]byte, error) {
	if n < 0xfd {
		return []byte{byte(n)}, nil
	}

	if n < 0x10000 {
		result := make([]byte, 3)
		result[0] = 0xfd
		binary.LittleEndian.PutUint16(result[1:], uint16(n))
		return result, nil
	}

	if n < 0x100000000 {
		result := make([]byte, 5)
		result[0] = 0xfe
		binary.LittleEndian.PutUint32(result[1:], uint32(n))
		return result, nil
	}

	result := make([]byte, 9)
	result[0] = 0xff
	binary.LittleEndian.PutUint64(result[1:], n)
	return result, nil
}
