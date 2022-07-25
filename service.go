package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	CSVSeparator = ';'

	csvID               = 0
	csvName             = 1
	csvAge              = 2
	csvBalanced         = 3
	csvPreviousBalanced = 4
	csvAverageBalanced  = 5
	csvFreeTransfer     = 6

	BonusBalance = 10
)

var NumberThreadFirst int32
var NumberThreadSecond int32
var NumberThreadThird int32

var wg sync.WaitGroup
var BankBudget float64 = 1000
var lCh chan *EODData
var mu sync.Mutex

func readBeforeEODCSV(fileName string) []*EODData {
	fullFilePath := fileName
	f, err := os.Open(fullFilePath)
	if err != nil {
		log.Fatal("Unable to read input file "+fullFilePath, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = CSVSeparator
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+fullFilePath, err)
	}

	return parseEODCSVRowToBeforeEODData(records)
}

func parseEODCSVRowToBeforeEODData(records [][]string) []*EODData {
	var list []*EODData
	for i, row := range records {
		// skip header data
		if i == 0 {
			continue
		}
		o := &EODData{}

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

func calculateAvargeBalance(q *EODData) {
	defer wg.Done()
	q.AverageBalanced = (q.Balanced + q.PreviousBalanced) / 2
	q.FirstThreadNumber = atomic.AddInt32(&NumberThreadFirst, 1)
}

func calculateBenefit(q *EODData) {
	defer wg.Done()
	if q.Balanced >= 100 && q.Balanced <= 150 {
		q.FreeTransfer = 5
	} else if q.Balanced > 150 {
		q.FreeTransfer = 5 + 25
	}
	q.SecondThreadNumber = atomic.AddInt32(&NumberThreadSecond, 1)
}

func calculateBonus() {
	defer wg.Done()
	noOfThread := atomic.AddInt32(&NumberThreadThird, 1)
	for q := range lCh {
		q.ThirdThreadNumber = noOfThread
		if q.ID >= 1 && q.ID <= 100 {
			mu.Lock()
			if BankBudget >= BonusBalance {
				q.Balanced += BonusBalance
				BankBudget -= BonusBalance
			}
			mu.Unlock()
		}
	}
}

func processEOD(EODList []*EODData) {

	for _, eod := range EODList {
		wg.Add(1)
		go calculateAvargeBalance(eod)
		wg.Add(1)
		go calculateBenefit(eod)
	}

	wg.Wait()

	lCh = make(chan *EODData)
	maxNumberOfChannel := 8
	for i := 0; i < maxNumberOfChannel; i++ {
		wg.Add(1)
		go calculateBonus()
	}
	for _, eod := range EODList {
		lCh <- eod
	}
	close(lCh)
	wg.Wait()
}

func writeAfterEODCSV(fileName string, records []*EODData) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = CSVSeparator

	var data [][]string
	// write header firest
	rowHeader := []string{"id", "Nama", "Age", "Balanced", "No 2b Thread-No", "No 3 Thread-No", "Previous Balanced", "Average Balanced",
		"No 1 Thread-No", "Free Transfer", "No 2a Thread-No"}
	data = append(data, rowHeader)

	// write csv body
	for _, record := range records {
		row := []string{strconv.Itoa(record.ID), record.Name, strconv.Itoa(record.Age), fmt.Sprintf("%.2f", record.Balanced), strconv.Itoa(int(record.SecondThreadNumber)),
			strconv.Itoa(int(record.ThirdThreadNumber)), fmt.Sprintf("%.2f", record.PreviousBalanced), fmt.Sprintf("%.2f", record.AverageBalanced), strconv.Itoa(int(record.FirstThreadNumber)),
			strconv.Itoa(record.FreeTransfer), strconv.Itoa(int(record.SecondThreadNumber))}
		data = append(data, row)
	}
	csvwriter.WriteAll(data)
}
