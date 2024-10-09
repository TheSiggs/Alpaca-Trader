package strategies

import (
	"encoding/json"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/nettis/alpaca-trader/config"
	"github.com/nettis/alpaca-trader/entities"
)


func DividendFlipper() {
	var client entities.TradingClient
	client.Config = config.Setup()
	client.Client = alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    client.Config.AlpacaConfig.APIKey,
		APISecret: client.Config.AlpacaConfig.APISecret,
		BaseURL:   client.Config.AlpacaConfig.BaseURL,
	})

	log.Println("Local Time Zone Location:", time.Local)

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
		Symbol:         SymbolToTrade,
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
