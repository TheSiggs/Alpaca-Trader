package main

import (
	"encoding/json"

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
	log.Printf("%+v\n\n", *acct)

	closed_orders, err := client.Client.CloseAllPositions(alpaca.CloseAllPositionsRequest{CancelOrders: true})
	if err != nil {
		log.Println(err)
	}
	jsonClosedOrders, err := json.Marshal(closed_orders)
	if err == nil {
		log.Printf("Positions closed: %+v\n", string(jsonClosedOrders))
	}

	// Look for a stock to buy this stock should have an ex dividend date for tomorrow in which we will sell it
	dividends := client.UpcomingDividends()
	if len(dividends) == 0 {
		log.Println("No stocks to buy today")
		return
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

	log.Printf("Trading On Symbol: %s\n\n", SymbolToTrade)
	// Buy the stock
	notional := acct.Cash // Amount in dollars you want to invest
	orderReq := alpaca.PlaceOrderRequest{
		Symbol: SymbolToTrade,
		Notional:       &notional,
		Side:           alpaca.Buy,       // Order side: Buy or Sell
		Type:           alpaca.Market,    // Order type: Market or Limit
		TimeInForce:    alpaca.Day,       // Time in force: Day, GTC, etc.
		PositionIntent: alpaca.BuyToOpen, // Opens a long position,
	}
	// Submit the order
	order, err := client.Client.PlaceOrder(orderReq) // FIXME: Error placing order: invalid position_intent specified (HTTP 422, Code 40010001)
	if err != nil {
		log.Println("Error placing order:", err)
		return
	}
	orderJSON, err := json.Marshal(order)
	if err == nil {
		log.Printf("Order placed for %s: %+v", SymbolToTrade, string(orderJSON))
	}
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
