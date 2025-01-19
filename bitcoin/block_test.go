package bitcoin_test

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func TestBlock(t *testing.T) {

	t.Run("ParseBlock return block with correct values", func(t *testing.T) {
		hexString := "000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if block.Version != 536870912 {
			t.Errorf("expected version 536870912, got %d", block.Version)
		}

		if block.PreviousBlock == nil {
			t.Errorf("expected previous block not to be nil")
		}

		if block.MerkleRoot == nil {
			t.Errorf("expected merkle root not to be nil")
		}

		if block.Timestamp != 1496586576 {
			t.Errorf("expected timestamp 1496586576, got %d", block.Timestamp)
		}

		bits := binary.LittleEndian.Uint32(block.Bits)
		if bits != 402774100 {
			t.Errorf("expected bits 402774100, got %d", bits)
		}

		if block.Nonce != 773550753 {
			t.Errorf("expected nonce 773550753, got %d", block.Nonce)
		}
	})

	t.Run("Serialize", func(t *testing.T) {
		hexString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		serialized, err := block.Serialize()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if hex.EncodeToString(serialized) != hexString {
			t.Errorf("expected serialized block to be %s, got %s", hexString, hex.EncodeToString(serialized))
		}
	})

	t.Run("Hash", func(t *testing.T) {
		hexString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		blockHash, err := block.Hash()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectedHash := "0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523"
		if hex.EncodeToString(blockHash) != expectedHash {
			t.Errorf("expected hash to be %x, got %x", expectedHash, blockHash)
		}
	})

	t.Run("Is BIP9 readiness", func(t *testing.T) {
		hexString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !block.BIP9() {
			t.Errorf("expected block to not be BIP9 ready")
		}
	})

	t.Run("Is not BIP9 readiness", func(t *testing.T) {
		hexString := "0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if block.BIP9() {
			t.Errorf("expected block to not be BIP9 ready")
		}
	})

	t.Run("Is BIP91 readiness", func(t *testing.T) {
		hexString := "1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !block.BIP91() {
			t.Errorf("expected block to be BIP91 ready")
		}
	})

	t.Run("Is not BIP91 readiness", func(t *testing.T) {
		hexString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if block.BIP91() {
			t.Errorf("expected block to not be BIP91 ready")
		}
	})

	t.Run("Is BIP141 readiness", func(t *testing.T) {
		hexString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !block.BIP141() {
			t.Errorf("expected block to be BIP141 ready")
		}
	})

	t.Run("Is not BIP141 readiness", func(t *testing.T) {
		hexString := "0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if block.BIP141() {
			t.Errorf("expected block to not be BIP141 ready")
		}
	})

	t.Run("Bits to target", func(t *testing.T) {
		hexString := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
		data, _ := hex.DecodeString(hexString)

		block, err := bitcoin.ParseBlock(data)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		target := block.Target()
		expectedTarget := new(big.Int)
		expectedTarget.SetString("13ce9000000000000000000000000000000000000000000", 16)
		if target.Cmp(expectedTarget) != 0 {
			t.Errorf("expected target to be %x, got %x", expectedTarget, target)
		}
	})

}
