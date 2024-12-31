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

func TestEncode(t *testing.T) {
	tests := []struct {
		input    *big.Int
		expected []byte
	}{
		{big.NewInt(1), []byte{0x01}},
		{big.NewInt(0xfc), []byte{0xfc}},
		{big.NewInt(0xfd), []byte{0xfd, 0xfd, 0x00}},
		{big.NewInt(0x10000), []byte{0xfe, 0x00, 0x00, 0x01, 0x00}},
		{big.NewInt(0x10000000), []byte{0xfe, 0x00, 0x00, 0x00, 0x10}},
	}

	for _, test := range tests {
		result, err := varint.Encode(test.input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bytes.Equal(result, test.expected) {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}
