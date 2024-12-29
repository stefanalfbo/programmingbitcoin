package endian_test

import (
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

func TestLittleEndianToBigInt(t *testing.T) {
	t.Run("LittleEndianToInt is valid", func(t *testing.T) {
		bytes := []byte{0x01, 0x00, 0x00}
		expected := big.NewInt(1)

		result := endian.LittleEndianToBigInt(bytes)
		if result.Cmp(expected) != 0 {
			t.Errorf("LittleEndianToBigInt: got %v, expected %v", result, expected)
		}
	})

	t.Run("LittleEndianToInt is valid", func(t *testing.T) {
		bytes := []byte{0x00, 0x01, 0x00}
		expected := big.NewInt(256)

		result := endian.LittleEndianToBigInt(bytes)
		if result.Cmp(expected) != 0 {
			t.Errorf("LittleEndianToBigInt: got %v, expected %v", result, expected)
		}
	})
}

func TestBigIntToLittleEndian(t *testing.T) {
	t.Run("BigIntToLittleEndian is valid", func(t *testing.T) {
		n := big.NewInt(1)
		length := 3
		expected := []byte{0x01, 0x00, 0x00}

		result := endian.BigIntToLittleEndian(n, length)
		if len(result) != length {
			t.Errorf("BigIntToLittleEndian: got %v, expected %v", result, expected)
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("BigIntToLittleEndian: got %v, expected %v", result, expected)
			}
		}
	})

	t.Run("BigIntToLittleEndian is valid", func(t *testing.T) {
		n := big.NewInt(256)
		length := 3
		expected := []byte{0x00, 0x01, 0x00}

		result := endian.BigIntToLittleEndian(n, length)
		if len(result) != length {
			t.Errorf("BigIntToLittleEndian: got %v, expected %v", result, expected)
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("BigIntToLittleEndian: got %v, expected %v", result, expected)
			}
		}
	})
}
