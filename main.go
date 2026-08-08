package main

import (
	"fmt"
	"sync"

	"github.com/wnjoon/go-yfinance/pkg/ticker"
	"github.com/wnjoon/go-yfinance/pkg/utils"
)

type tickerInfo struct {
	name        string
	symbol      string
	marketPrice float64
}

func processTickers(tickers []string) []tickerInfo {
	var wg sync.WaitGroup
	numTickers := len(tickers)
	res := make([]tickerInfo, numTickers)

	// use go routine to fetch all tickers in parallel
	for i := range len(tickers) {

		// .Go automatically calls Add / Done
		wg.Go(
			func() {
				ti := getTickerInfo(tickers[i])

				// Since each goroutine maps to a specifc index, don't need a mutex for updation, routines will ONLY update their designated spot in the slice.
				res[i] = ti
			},
		)
	}

	// wait for all go routines to finish
	wg.Wait()
	return res
}

func getTickerInfo(tickerStr string) tickerInfo {
	formattedTicker := utils.FormatYahooTicker(tickerStr, "XTSE")
	t, err := ticker.New(formattedTicker)

	// TODO: deal with errors better somehow
	// possibly return better string, add error code to tickerInfo
	if err != nil {
		fmt.Println("uh oh", err)
	}

	quote, err := t.Quote()
	if err != nil {
		fmt.Println("uh oh 2", err)
	}

	res := tickerInfo{
		name:        quote.LongName,
		symbol:      quote.Symbol,
		marketPrice: quote.RegularMarketPrice,
	}

	return res
}

func main() {
	// Plan for this:

	// Part 1: Naive fetching
	// Take in list of tickers - DONE
	// Fetch them in parallel - DONE
	// Wait for them all, then print them all together when they all arrive - DONE

	// Part 2: error handling
	// what do i do if all of them don't arrive?
	// timeout?
	// what if it comes and it is bad data? API fails

	// Part 3: Modularize
	// should not be in the main pkg, main is for coordination

	// Part 4: Doc
	// what did u learn, where do these logs go in docs struct??

	tickers := [3]string{"TD", "RY", "VFV"}
	for {
		infoState := processTickers(tickers[:])
		fmt.Println(infoState)
	}

}
