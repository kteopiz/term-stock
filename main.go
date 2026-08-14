package main

import (
	"fmt"
	"sync"

	"github.com/wnjoon/go-yfinance/pkg/ticker"
	"github.com/wnjoon/go-yfinance/pkg/utils"
)

// TODO: Move all to proper module (tickerInfo smth)
type tickerInfo struct {
	name        string
	symbol      string
	marketPrice float64
	isValid     bool
}

var EmptyTickerInfo = tickerInfo{
	"",
	"",
	-1.0,
	false,
}

func NewTickerInfo(name string, symbol string, marketPrice float64) tickerInfo {
	new := tickerInfo{
		name:        name,
		symbol:      symbol,
		marketPrice: marketPrice,
	}

	// Want to flag negative MP since implies error from fetch
	// Since MP input is out of control from yfinance set based on value fetched instead of setting manually via. tickerInfo{}
	new.isValid = marketPrice >= 0

	return new
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
				ti, err := getTickerInfo(tickers[i])
				if err != nil {
					fmt.Printf("error: %s\n", err)
					// flag isValid to show clear error and enforce action on TUI (TBD on action(s) to take: red line, disappearing ticker, show error in tabular section, etc.)
					res[i].isValid = false
				} else {
					// Since each goroutine maps to a specifc index, don't need a mutex for updation, routines will ONLY update their designated spot in the slice.
					res[i] = ti
				}
			},
		)
	}

	// wait for all go routines to finish
	wg.Wait()
	return res
}

// Error Handling
// Tickers will come in as trusted. Trusted defined as a valid symbol, non-empty symbol which exists on the TSX
// FUTURE TODO: Other functions will parse .csv portfolio files or take in user input for new tickers and clean them such that they come in trusted
func getTickerInfo(tickerStr string) (tickerInfo, error) {

	// TODO: Remove in the future... user input should be cleaned in future to not require this
	// Edge Case: Must check if empty ticker before appending .TO suffix, since ".TO" is a valid ticker.
	if tickerStr == "" {
		return EmptyTickerInfo, fmt.Errorf("empty ticker provided")
	}

	// add .TO yfinance suffix
	formattedTicker := utils.FormatYahooTicker(tickerStr, "XTSE")

	t, err := ticker.New(formattedTicker)
	if err != nil {
		return EmptyTickerInfo, fmt.Errorf("failed to create new ticker: %w", err)
	}

	// ERROR SOURCES

	// https://github.com/wnjoon/go-yfinance/blob/main/pkg/client/auth.go
	// getWithCrumb -> addCrumbToParams -> GetCrumb ->  (if no crumb) RefreshAuth -> fetchBasic or fetchCSRF
	// GetCrumb -> returns crumb, if DNE does fetch for it
	// RefreshAuth -> occurs on above case of DNE crumb, returns a crumb fetched in one of two ways: fetchBasic or fetchCSRF two diff strats
	// fetchBasic -> cookie or crumb failure fully fails on the client.Get all rely on if the GET requests is successful or not
	// fetchBasic(2) ->  429 rate limit, and invalid response comes directly from the response body
	// fetchCSRF -> same as fetchBasic reliance on GET/POST requests, then 429 or invalid body

	// errors from getWithCrumb boil down to 3 cases then at the core:
	// 1. some GET or POST request failure
	// 2. 429 rate limit error
	// https://github.com/wnjoon/go-yfinance/blob/main/pkg/client/auth.go#L299
	// 3a. invalid body response from fetched response for crumb in fetchCSRF/fetchBasic

	// Rest of errors come from Quote code post parse/unmarshal of JSON response
	// https://github.com/wnjoon/go-yfinance/blob/main/pkg/ticker/quote.go
	// 3b. invalid body response from WrapNotFoundError defined error in Quote code
	// 4. API error from above Quote code as well

	// Creating a quote creates a fetch to yfinance, uses crumb auth
	quote, err := t.Quote()
	if err != nil {
		return EmptyTickerInfo, fmt.Errorf("failed to get quote for %s %w", formattedTicker, err)
	}

	res := NewTickerInfo(
		quote.LongName,
		quote.Symbol,
		quote.RegularMarketPrice,
	)

	// Edge: yfinance returns a negative stock price
	if !res.isValid {
		return EmptyTickerInfo, fmt.Errorf("invalid negative market price for %s", res.symbol)
	}

	return res, nil
}

func main() {
	// Plan for this:

	// Part 1: Naive fetching
	// Take in list of tickers - DONE
	// Fetch them in parallel - DONE
	// Wait for them all, then print them all together when they all arrive - DONE

	// Part 2: error handling
	// Find all unique error sources & ensure they can bubble up to future controller for action - DONE, covers bad/failed requests and rate limiting errors
	// Pass an error on timeout

	// Part 3: Modularize
	// should not be in the main pkg, main is for coordination

	// Part 4: Doc
	// what did u learn, where do these logs go in docs struct??

	tickers := [4]string{"TD", "RY", "VFV", "DASKDKWAD"}
	for range 5 {
		infoState := processTickers(tickers[:])
		fmt.Println(infoState)
	}
}
