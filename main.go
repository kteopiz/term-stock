package main

import (
	"fmt"

	"github.com/wnjoon/go-yfinance/pkg/ticker"
	"github.com/wnjoon/go-yfinance/pkg/utils"
)

func main() {
	rawTicker := "H"
	formattedTicker := utils.FormatYahooTicker(rawTicker, "XTSE")
	t, err := ticker.New(formattedTicker)
	if err != nil {
		fmt.Println("uh oh", err)
	}

	quote, err := t.Quote()
	if err != nil {
		fmt.Println("uh oh 2", err)
	}

	fmt.Printf(" Name: %s, Symbol: %s\n Exchange: %s\n Market Price: %.2f\n Currency: %s\n Market State: %s",
		quote.LongName,
		quote.Symbol,
		quote.ExchangeName,
		quote.RegularMarketPrice,
		quote.Currency,
		quote.MarketState,
	)

}
