package utils

import (
	"encoding/csv"
	"os"
)

// The ReadCsv function reads a CSV file from the specified path and returns its contents as a 2D slice
// of strings along with any encountered errors.
func ReadCsv(path string) ([][]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	} else {
		defer file.Close()
	}

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()

	if err != nil {
		return nil, err
	}

	return records, nil
}
