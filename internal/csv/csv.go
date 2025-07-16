package csv

import "fmt"

type CSV struct {
	Header []string
	Body   [][]string
}

func NewCSV(records [][]string) (*CSV, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("empty csv")
	}
	return &CSV{
		Header: records[0],
		Body:   records[1:],
	}, nil
}
