package main

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	Config "github.com/nettis/dividend-trader-go/config"

	"context"
	"log"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

type TradingClient struct {
	Client *alpaca.Client
	Config Config.Config
}

func main() {
	var client TradingClient
    client.Config = Config.Setup() 
	client.Client = alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    client.Config.AlpacaConfig.APIKey,
		APISecret: client.Config.AlpacaConfig.APISecret,
		BaseURL:   client.Config.AlpacaConfig.BaseURL,
	})

	acct, err := client.Client.GetAccount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n\n", *acct)

	_, err = client.Client.CloseAllPositions(alpaca.CloseAllPositionsRequest{CancelOrders: true})
	if err != nil {
		fmt.Println(err)
	}

	// Look for a stock to buy this stock should have an ex dividend date for tomorrow in which we will sell it
	dividends := client.UpcomingDividends()
	if len(dividends) == 0 {
		fmt.Println("No stocks to buy today")
		return
	}
	var SymbolToTrade string
	for _, dividend := range dividends {
		fmt.Printf("Symbol: %s, ExDate: %s, DivAmt: %.5f\n", dividend.Ticker, dividend.ExDividendDate, dividend.CashAmount)
		if client.CheckSymbol(dividend.Ticker) {
			SymbolToTrade = dividend.Ticker
			break
		}
	}

	fmt.Printf("Trading On Symbol: %s\n\n", SymbolToTrade)
    // Buy the stock 
	notional := acct.Cash // Amount in dollars you want to invest
	orderReq := alpaca.PlaceOrderRequest{
		Symbol:      SymbolToTrade,
		Notional:    &notional,
		Side:        alpaca.Buy,    // Order side: Buy or Sell
		Type:        alpaca.Market, // Order type: Market or Limit
		TimeInForce: alpaca.Day,    // Time in force: Day, GTC, etc.
	}
	// Submit the order
	_, err = client.Client.PlaceOrder(orderReq) // FIXME: Error placing order: invalid position_intent specified (HTTP 422, Code 40010001)
	if err != nil {
		fmt.Println("Error placing order:", err)
	}
}

func (client *TradingClient) CheckSymbol(symbol string) bool {
	asset, err := client.Client.GetAsset(symbol)
	if err != nil {
		fmt.Println(err)
		return false
	}
    if asset.Status == alpaca.AssetInactive {
        fmt.Printf("%s is inactive and cannot be traded currently", symbol) 
        return false
    }
	fmt.Println(asset)
	return true
}

func (client *TradingClient) UpcomingDividends() []models.Dividend {
	var dividends []models.Dividend

	c := polygon.New(client.Config.PolygonConfig.APIKey)

	currentDate := time.Now()
	currentDate = currentDate.AddDate(0, 0, 1)
	year, month, day := currentDate.Date()
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
