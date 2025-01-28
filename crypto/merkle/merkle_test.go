package merkle_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/crypto/merkle"
)

func TestParent(t *testing.T) {
	expectedParent := "8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd"
	left, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	right, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")

	parent := merkle.Parent(left, right)

	parentHex := hex.EncodeToString(parent)
	if parentHex != expectedParent {
		t.Errorf("Parent was incorrect, got: %s, want: %s.", parentHex, expectedParent)
	}
}

func TestParentLevel(t *testing.T) {
	expectedParentHashes := []string{
		"8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd",
		"7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800",
		"3ecf6115380c77e8aae56660f5634982ee897351ba906a6837d15ebc3a225df0",
	}
	hashes := make([][]byte, 5)
	hashes[0], _ = hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	hashes[1], _ = hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	hashes[2], _ = hex.DecodeString("f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0")
	hashes[3], _ = hex.DecodeString("3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181")
	hashes[4], _ = hex.DecodeString("10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae")

	parentHashes := merkle.ParentLevel(hashes)

	for i, parent := range parentHashes {
		parentHex := hex.EncodeToString(parent)
		if parentHex != expectedParentHashes[i] {
			t.Errorf("ParentLevel was incorrect, got: %s, want: %s.", parentHex, expectedParentHashes[i])
		}
	}
}

func TestRoot(t *testing.T) {
	expectedRoot := "acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6"
	hashes := make([][]byte, 12)
	hashes[0], _ = hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	hashes[1], _ = hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	hashes[2], _ = hex.DecodeString("f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0")
	hashes[3], _ = hex.DecodeString("3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181")
	hashes[4], _ = hex.DecodeString("10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae")
	hashes[5], _ = hex.DecodeString("7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161")
	hashes[6], _ = hex.DecodeString("8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc")
	hashes[7], _ = hex.DecodeString("dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877")
	hashes[8], _ = hex.DecodeString("b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59")
	hashes[9], _ = hex.DecodeString("95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c")
	hashes[10], _ = hex.DecodeString("2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908")
	hashes[11], _ = hex.DecodeString("b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0")

	root := merkle.Root(hashes)

	rootHex := hex.EncodeToString(root)
	if rootHex != expectedRoot {
		t.Errorf("Root was incorrect, got: %s, want: %s.", rootHex, expectedRoot)
	}
}

func TestNewMerkleTree(t *testing.T) {
	tree := merkle.NewMerkleTree(9)

	if len(tree.Nodes[0]) != 1 {
		t.Errorf("NewMerkleTree was incorrect, got: %d, want: %d.", len(tree.Nodes[0]), 1)
	}

	if len(tree.Nodes[1]) != 2 {
		t.Errorf("NewMerkleTree was incorrect, got: %d, want: %d.", len(tree.Nodes[1]), 2)
	}

	if len(tree.Nodes[2]) != 3 {
		t.Errorf("NewMerkleTree was incorrect, got: %d, want: %d.", len(tree.Nodes[2]), 3)
	}

	if len(tree.Nodes[3]) != 5 {
		t.Errorf("NewMerkleTree was incorrect, got: %d, want: %d.", len(tree.Nodes[3]), 5)
	}

	if len(tree.Nodes[4]) != 9 {
		t.Errorf("NewMerkleTree was incorrect, got: %d, want: %d.", len(tree.Nodes[4]), 9)
	}
}

func TestString(t *testing.T) {
	hashes := make([][]byte, 16)
	hashes[0], _ = hex.DecodeString("9745f7173ef14ee4155722d1cbf13304339fd00d900b759c6f9d58579b5765fb")
	hashes[1], _ = hex.DecodeString("5573c8ede34936c29cdfdfe743f7f5fdfbd4f54ba0705259e62f39917065cb9b")
	hashes[2], _ = hex.DecodeString("82a02ecbb6623b4274dfcab82b336dc017a27136e08521091e443e62582e8f05")
	hashes[3], _ = hex.DecodeString("507ccae5ed9b340363a0e6d765af148be9cb1c8766ccc922f83e4ae681658308")
	hashes[4], _ = hex.DecodeString("a7a4aec28e7162e1e9ef33dfa30f0bc0526e6cf4b11a576f6c5de58593898330")
	hashes[5], _ = hex.DecodeString("bb6267664bd833fd9fc82582853ab144fece26b7a8a5bf328f8a059445b59add")
	hashes[6], _ = hex.DecodeString("ea6d7ac1ee77fbacee58fc717b990c4fcccf1b19af43103c090f601677fd8836")
	hashes[7], _ = hex.DecodeString("457743861de496c429912558a106b810b0507975a49773228aa788df40730d41")
	hashes[8], _ = hex.DecodeString("7688029288efc9e9a0011c960a6ed9e5466581abf3e3a6c26ee317461add619a")
	hashes[9], _ = hex.DecodeString("b1ae7f15836cb2286cdd4e2c37bf9bb7da0a2846d06867a429f654b2e7f383c9")
	hashes[10], _ = hex.DecodeString("9b74f89fa3f93e71ff2c241f32945d877281a6a50a6bf94adac002980aafe5ab")
	hashes[11], _ = hex.DecodeString("b3a92b5b255019bdaf754875633c2de9fec2ab03e6b8ce669d07cb5b18804638")
	hashes[12], _ = hex.DecodeString("b5c0b915312b9bdaedd2b86aa2d0f8feffc73a2d37668fd9010179261e25e263")
	hashes[13], _ = hex.DecodeString("c9d52c5cb1e557b92c84c52e7c4bfbce859408bedffc8a5560fd6e35e10b8800")
	hashes[14], _ = hex.DecodeString("c555bc5fc3bc096df0a0c9532f07640bfb76bfe4fc1ace214b8b228a1297a4c2")
	hashes[15], _ = hex.DecodeString("f9dbfafc3af3400954975da24eb325e326960a25b87fffe23eef3e7ed2fb610e")

	tree := merkle.NewMerkleTree(len(hashes))
	tree.Nodes[4] = make([][]byte, len(hashes))

	copy(tree.Nodes[4], hashes)

	tree.Nodes[3] = merkle.ParentLevel(tree.Nodes[4])
	tree.Nodes[2] = merkle.ParentLevel(tree.Nodes[3])
	tree.Nodes[1] = merkle.ParentLevel(tree.Nodes[2])
	tree.Nodes[0] = merkle.ParentLevel(tree.Nodes[1])

	expected := "depth 0: *597c4baf*"

	result := tree.String()

	if strings.Contains(result, expected) == false {
		t.Errorf("String was incorrect, got: %s, want: %s.", result, expected)
	}
}
