package hash_test

import (
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/hash"
)

func TestDoubleSHA256(t *testing.T) {
	expected := "9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50"

	hash := hash.Hash256([]byte("hello"))

	if hex.EncodeToString(hash) != expected {
		t.Errorf("Expected %s but got %s", expected, hex.EncodeToString(hash))
	}
}

func TestSHA256RIPEMD160(t *testing.T) {
	expected := "b6a9c8c230722b7c748331a8b450f05566dc7d0f"

	hash := hash.Hash160([]byte("hello"))

	if hex.EncodeToString(hash) != expected {
		t.Errorf("Expected %s but got %s", expected, hex.EncodeToString(hash))
	}
}

func TestSHA1(t *testing.T) {
	expected := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"

	hash := hash.HashSHA1([]byte("hello"))

	if hex.EncodeToString(hash) != expected {
		t.Errorf("Expected %s but got %s", expected, hex.EncodeToString(hash))
	}
}
