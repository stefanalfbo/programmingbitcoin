package base58_test

import (
	"encoding/hex"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/encoding/base58"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", "9MA8fRQrT4u8Zj8ZRd6MAiiyaxb2Y1CMpvVkHQu5hVM6"},
		{"eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", "4fE3H2E6XMp4SsxtwinF7w9a34ooUrwWe4WsW1458Pd"},
		{"c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", "EQJsjkd6JaGwxrjEhfeqPenqHwrBmPQZjJGNSCHBkcF7"},
	}

	for _, test := range tests {
		inputBytes, _ := hex.DecodeString(test.input)
		encoded := base58.Encode(inputBytes)
		if encoded != test.expected {
			t.Errorf("Encode(%s): got %s, expected %s", test.input, encoded, test.expected)
		}
	}
}

func TestDecode(t *testing.T) {
	expected := "507b27411ccf7f16f10297de6cef3f291623eddf"
	address := "mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf"

	decoded, err := base58.Decode(address)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	decodedAsHex := hex.EncodeToString(decoded)

	if decodedAsHex != expected {
		t.Errorf("Decode(%s): got %s, expected %s", address, decodedAsHex, expected)
	}
}
