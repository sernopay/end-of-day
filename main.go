package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("file name is required")
	}
	previousEODList := ReadBeforeEODCSV(os.Args[1])
	fmt.Println(previousEODList)
}
