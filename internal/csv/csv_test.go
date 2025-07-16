package csv

import (
	"reflect"
	"testing"
)

func TestNewCSV(t *testing.T) {
	records := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
		{"value3", "value4"},
	}
	csv, err := NewCSV(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedHeader := []string{"header1", "header2"}
	if !reflect.DeepEqual(csv.Header, expectedHeader) {
		t.Errorf("expected header %v, but got %v", expectedHeader, csv.Header)
	}
	expectedBody := [][]string{
		{"value1", "value2"},
		{"value3", "value4"},
	}
	if !reflect.DeepEqual(csv.Body, expectedBody) {
		t.Errorf("expected body %v, but got %v", expectedBody, csv.Body)
	}
}

func TestNewCSV_empty(t *testing.T) {
	_, err := NewCSV([][]string{})
	if err == nil {
		t.Fatal("expected error, but got nil")
	}
}
