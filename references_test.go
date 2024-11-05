//nolint:gochecknoglobals,gochecknoinits // test data
package uid_test

import (
	"embed"
	"encoding/csv"
	"io"
)

//go:embed samples.csv
var files embed.FS

var sampleData []sampleid

type sampleid struct{ Canonical, B32, B64 string }

func init() {
	csvFile, err := files.Open("samples.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(csvFile)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		entry := sampleid{row[1], row[2], row[4]}
		sampleData = append(sampleData, entry)
	}
}
