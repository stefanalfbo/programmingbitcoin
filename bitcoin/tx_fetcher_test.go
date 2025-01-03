package bitcoin_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stefanalfbo/programmingbitcoin/bitcoin"
)

func MockFetcher(txId string, isTestnet bool) ([]byte, error) {
	path := filepath.Join("testdata", "15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521.raw")

	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return content, nil
}

func TestTxFetcher(t *testing.T) {
	t.Skip("Implementation is not done yet.")
	isTestnet := false
	isFresh := false
	txFetcher := bitcoin.NewTxFetcher(MockFetcher, isTestnet)
	txId := "15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521"

	tx, err := txFetcher.Fetch(txId, isFresh)
	if err != nil {

		t.Fatalf("error fetching transaction: %v", err)
	}

	fmt.Printf("tx: %v\n", tx)
}
