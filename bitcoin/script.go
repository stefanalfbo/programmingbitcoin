package bitcoin

import (
	"io"

	"github.com/stefanalfbo/programmingbitcoin/encoding/varint"
)

func ParseScript(data io.Reader) ([]byte, error) {
	length, err := varint.Decode(data)
	if err != nil {
		return nil, err
	}

	scriptLength := int(length.Int64())

	for i := 0; i < scriptLength; i++ {
		current := make([]byte, 1)
		_, err := data.Read(current)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
