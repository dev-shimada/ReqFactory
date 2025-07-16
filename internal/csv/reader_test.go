package csv

import (
	"bytes"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	// BOM付きのCSVデータ
	bom := []byte{0xef, 0xbb, 0xbf}
	csvData := []byte(`header1,header2
value1,value2
`)
	dataWithBom := append(bom, csvData...)

	reader := bytes.NewReader(dataWithBom)
	records, err := Read(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}

	if !reflect.DeepEqual(records, expected) {
		t.Errorf("expected %v, but got %v", expected, records)
	}
}
