package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	FilePath        = "before-eod/"
	CSVSeparatedKey = ';'

	csvID               = 0
	csvName             = 1
	csvAge              = 2
	csvBalanced         = 3
	csvPreviousBalanced = 4
	csvAverageBalanced  = 5
	csvFreeTransfer     = 6
)

func ReadBeforeEODCSV(fileName string) []*BeforeEODData {
	fullFilePath := FilePath + fileName
	f, err := os.Open(fullFilePath)
	if err != nil {
		log.Fatal("Unable to read input file "+fullFilePath, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = CSVSeparatedKey
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+fullFilePath, err)
	}

	return parseEODCSVRowToBeforeEODData(records)
}

func parseEODCSVRowToBeforeEODData(records [][]string) []*BeforeEODData {
	var list []*BeforeEODData
	for i, row := range records {
		// skip header data
		if i == 0 {
			continue
		}
		o := &BeforeEODData{}

		// id
		if id, err := strconv.Atoi(row[csvID]); err != nil {
			log.Fatal(fmt.Sprintf("failed to parse ID %s in row %d, err : %s", row[csvID], i, err))
		} else {
			o.ID = id
		}

		// name
		o.Name = row[csvName]

		// age
		if age, err := strconv.Atoi(row[csvAge]); err != nil {
			log.Fatal(fmt.Sprintf("failed to parse Age %s in row %d, err : %s", row[csvAge], i, err))
		} else {
			o.Age = age
		}

		// balanced
		if balanced, err := strconv.ParseFloat(row[csvBalanced], 64); err != nil {
			log.Fatal(fmt.Sprintf("failed to parse Balanced %s in row %d, err : %s", row[csvBalanced], i, err))
		} else {
			o.Balanced = balanced
		}

		// previous balanced
		if prevBalanced, err := strconv.ParseFloat(row[csvPreviousBalanced], 64); err != nil {
			log.Fatal(fmt.Sprintf("failed to parse Previous Balanced %s in row %d, err : %s", row[csvPreviousBalanced], i, err))
		} else {
			o.PreviousBalanced = prevBalanced
		}

		// average balanced
		if avgBalanced, err := strconv.ParseFloat(row[csvAverageBalanced], 64); err != nil {
			log.Fatal(fmt.Sprintf("failed to parse average Balanced %s in row %d, err : %s", row[csvAverageBalanced], i, err))
		} else {
			o.AverageBalanced = avgBalanced
		}

		// free transfer
		if freeTransfer, err := strconv.Atoi(row[csvFreeTransfer]); err != nil {
			log.Fatal(fmt.Sprintf("failed to parse free transfer %s in row %d, err : %s", row[freeTransfer], i, err))
		} else {
			o.FreeTransfer = freeTransfer
		}

		list = append(list, o)
	}
	return list
}
