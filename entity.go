package main

type EODData struct {
	ID               int
	Name             string
	Age              int
	Balanced         float64
	PreviousBalanced float64
	AverageBalanced  float64
	FreeTransfer     int

	FirstThreadNumber  int32
	SecondThreadNumber int32
	ThirdThreadNumber  int32
}
