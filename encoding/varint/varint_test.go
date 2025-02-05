package varint_test

import (
	"bytes"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected uint64
	}{
		{[]byte{0x01}, 1},
		{[]byte{0xfd, 0x01, 0x00}, 1},
		{[]byte{0xfe, 0x01, 0x00, 0x00, 0x00}, 1},
		{[]byte{0xff, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 1},
	}

	for _, test := range tests {
		stream := bytes.NewReader(test.input)
		result, err := varint.Decode(stream)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != test.expected {
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
		input    uint64
		expected []byte
	}{
		{1, []byte{0x01}},
		{0xfc, []byte{0xfc}},
		{0xfd, []byte{0xfd, 0xfd, 0x00}},
		{0x10000, []byte{0xfe, 0x00, 0x00, 0x01, 0x00}},
		{0x10000000, []byte{0xfe, 0x00, 0x00, 0x00, 0x10}},
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
