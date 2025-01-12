package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/stefanalfbo/programmingbitcoin/encoding/endian"
)

type Fetcher func(txId string, isTestnet bool) ([]byte, error)

type TxFetcher struct {
	cache     map[string]*Tx
	fetcher   Fetcher
	isTestnet bool
}

func NewTxFetcher(fetcher Fetcher, isTestnet bool) *TxFetcher {
	return &TxFetcher{
		cache:     make(map[string]*Tx),
		fetcher:   fetcher,
		isTestnet: isTestnet,
	}
}

func getUrl(isTestnet bool) string {
	if isTestnet {
		return "https://mempool.space/testnet4"
	}
	return "https://mempool.space"
}

func MemPoolFetcher(txId string, isTestnet bool) ([]byte, error) {
	url := fmt.Sprintf("%s/api/tx/%s/raw", getUrl(isTestnet), txId)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching transaction: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	raw, err := hex.DecodeString(string(body))
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (txf *TxFetcher) Fetch(txId string, isFresh bool) (*Tx, error) {
	if isFresh || txf.cache[txId] == nil {
		raw, err := txf.fetcher(txId, txf.isTestnet)
		if err != nil {
			return nil, err
		}

		var tx *Tx
		if raw[4] == 0 {
			raw = append(raw[:4], raw[6:]...)
			tx, err = Parse(bytes.NewReader(raw), txf.isTestnet)
			tx.LockTime = endian.LittleEndianToInt32(raw[len(raw)-4:])
			if err != nil {
				return nil, err
			}
		} else {
			tx, err = Parse(bytes.NewReader(raw), txf.isTestnet)
			if err != nil {
				return nil, err
			}
		}

		if tx.Id() != txId {
			return nil, fmt.Errorf("not the same id: %s vs %s", tx.Id(), txId)
		}

		txf.cache[txId] = tx

	}
	// txf.cache[txId].isTestnet = isTestnet

	return txf.cache[txId], nil
}
