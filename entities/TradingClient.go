package entities

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/nettis/alpaca-trader/config"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)


type TradingClient struct {
	Client *alpaca.Client
	Config config.Config
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
	return true
}

func (client *TradingClient) UpcomingDividends() []models.Dividend {
	var dividends []models.Dividend

	c := polygon.New(client.Config.PolygonConfig.APIKey)

    // Fallback Date to guess next market open
	currentDate := time.Now()
	currentDate = currentDate.AddDate(0, 0, 1)
	year, month, day := currentDate.Date()

    // Use API to get next market open
	clock, err := client.Client.GetClock()
	if err == nil {
        year, month, day = clock.NextOpen.Date()
	}

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

func (client *TradingClient) TodaysDividends() []models.Dividend {
	var dividends []models.Dividend

	c := polygon.New(client.Config.PolygonConfig.APIKey)

    // Fallback Date to guess next market open
	currentDate := time.Now()
	year, month, day := currentDate.Date()

    // Use API to get next market open
	clock, err := client.Client.GetClock()
	if err != nil {
       log.Fatal(err) 
        return nil 
	}
    
    if !clock.IsOpen {
        log.Fatalln("Market is closed")
    }

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
