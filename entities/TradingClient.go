package entities

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nettis/alpaca-trader/config"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)


type TradingClient struct {
	Client *alpaca.Client
    MarketClient *marketdata.Client
	Config config.Config
}

type TradingClientInterface interface {
	CurrentDateTime(days int) (int, time.Month, int)
}

func (client *TradingClient) CheckSymbol(symbol string) bool {
	asset, err := client.Client.GetAsset(symbol)
	if err != nil {
		log.Println(err)
		return false
	}
	assetJSON, err := json.Marshal(asset)
	if err == nil {
		log.Printf("Symbol Found: %v", string(assetJSON))
	}
	if asset.Status == alpaca.AssetInactive {
		log.Printf("%s is inactive and cannot be traded currently", symbol)
		return false
	}
	if asset.Tradable == false {
		log.Printf("%s is currently not tradable", symbol)
		return false
	}
	if asset.Fractionable == false {
		log.Printf("%s if not fraction tradable", symbol)
		return false
	}
    if asset.Shortable == false {
		log.Printf("%s if not shortable", symbol)
		return false
    }
	return true
}

func (client *TradingClient) UpcomingDividends() []models.Dividend {
	var dividends []models.Dividend

	c := polygon.New(client.Config.PolygonConfig.APIKey)

    year, month, day := client.CurrentDate(1) 

	params := models.ListDividendsParams{}.
		WithExDividendDate(models.EQ, models.Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))).
		WithSort("cash_amount").
		WithLimit(100)

	iter := c.ListDividends(context.TODO(), params)
	for iter.Next() {
		dividends = append(dividends, iter.Item())
	}
	if iter.Err() != nil {
		log.Print(iter.Err())
	}

	return dividends
}


func (client *TradingClient) Dividends(year int, month time.Month, day int) []models.Dividend {
	var dividends []models.Dividend

	c := polygon.New(client.Config.PolygonConfig.APIKey)

	params := models.ListDividendsParams{}.
		WithExDividendDate(models.EQ, models.Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))).
		WithSort("cash_amount").
		WithLimit(100)

	iter := c.ListDividends(context.TODO(), params)
	for iter.Next() {
		dividends = append(dividends, iter.Item())
	}
	if iter.Err() != nil {
		log.Print(iter.Err())
	}

	return dividends
}


func (client *TradingClient) CurrentDate(days int) (int, time.Month, int) {
    // Fallback Date to guess next market open
	currentDate := time.Now()
    if days != 0 {
        currentDate = currentDate.AddDate(0, 0, days)
    } 
	year, month, day := currentDate.Date()

    client.CheckMarket()

    return year, month, day
}

func (client *TradingClient) CheckMarket() {
    // Use API to get next market open
	clock, err := client.Client.GetClock()
	if err != nil {
        log.Fatalf("Failed to get clock: %v", err)
	}
    if !clock.IsOpen {
        log.Fatalln("Market is closed")
    }
}

func (client *TradingClient) LargestDividendStock(year int, month time.Month, day int) (string, error) {
	// Look for a stock to buy this stock should have an ex dividend date for tomorrow in which we will sell it
	dividends := client.Dividends(year, month, day)
	if len(dividends) == 0 {
		log.Println("No stocks to buy today")
		return "", errors.New("No stocks to buy today")
	}

	var SymbolToTrade string
	for _, dividend := range dividends {
		dividendJSON, _ := json.Marshal(dividend)
		log.Printf("Checking Symbol: %v", string(dividendJSON))
		if client.CheckSymbol(dividend.Ticker) {
			SymbolToTrade = dividend.Ticker
			break
		}
	}
    return SymbolToTrade, nil
}
