package csv

import (
	"bytes"
	"encoding/csv"
	"io"
)

var bom = []byte{0xef, 0xbb, 0xbf}

func Read(r io.Reader) ([][]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimPrefix(data, bom)
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
