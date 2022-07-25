package main

func main() {
	inputFileName := "before-eod/Before Eod.csv"
	outputFileName := "after-eod/After Eod.csv"

	EODList := readBeforeEODCSV(inputFileName)
	processEOD(EODList)
	writeAfterEODCSV(outputFileName, EODList)
}
