package varint_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected *big.Int
	}{
		{[]byte{0x01}, big.NewInt(1)},
		{[]byte{0xfd, 0x01, 0x00}, big.NewInt(1)},
		{[]byte{0xfe, 0x01, 0x00, 0x00, 0x00}, big.NewInt(1)},
		{[]byte{0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, big.NewInt(1)},
	}

	for _, test := range tests {
		stream := bytes.NewReader(test.input)
		result, err := varint.Decode(stream)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Cmp(test.expected) != 0 {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}

func TestDecode_Error(t *testing.T) {
	tests := [][]byte{
		{0xfd},
		{0xfe},
		{0xff},
	}

	for _, test := range tests {
		stream := bytes.NewReader(test)
		_, err := varint.Decode(stream)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	}
}
