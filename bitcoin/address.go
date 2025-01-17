package bitcoin

import (
	"github.com/stefanalfbo/programmingbitcoin/encoding/base58"
)

func H160ToP2SHAddress(h160 []byte, isTestnet bool) string {
	prefix := byte(0x05)
	if isTestnet {
		prefix = 0xc4
	}

	return base58.Checksum(append([]byte{prefix}, h160...))
}

func H160ToP2PKHAddress(h160 []byte, isTestnet bool) string {
	prefix := byte(0x00)
	if isTestnet {
		prefix = 0x6f
	}

	return base58.Checksum(append([]byte{prefix}, h160...))
}
