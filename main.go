package main

import (
	"fmt"

	"term-stock/tickerinfo"
)

func main() {
	tickers := [3]string{"TD", "RY", "VFV"}
	for range 5 {
		infoState := tickerinfo.ProcessTickers(tickers[:])
		fmt.Println(infoState)
	}
}
