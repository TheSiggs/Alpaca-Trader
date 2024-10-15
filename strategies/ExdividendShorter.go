package strategies

import (
	"encoding/json"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nettis/alpaca-trader/config"
	"github.com/nettis/alpaca-trader/entities"
	"github.com/shopspring/decimal"
)

func ExdividendShorter() {
	var client entities.TradingClient
	client.Config = config.Setup()
	client.Client = alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    client.Config.AlpacaConfig.APIKey,
		APISecret: client.Config.AlpacaConfig.APISecret,
		BaseURL:   client.Config.AlpacaConfig.BaseURL,
	})

	client.MarketClient = marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    client.Config.AlpacaConfig.APIKey,
		APISecret: client.Config.AlpacaConfig.APISecret,
		BaseURL: client.Config.AlpacaConfig.MarketBaseURL,
	})

	log.Println("Local Time Zone Location:", time.Local)

	acct, err := client.Client.GetAccount()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n\n", *acct)

	// Look for a stock to buy this stock should have an ex dividend date for tomorrow in which we will sell it
	dividends := client.TodaysDividends()
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

	quote, err := client.MarketClient.GetLatestQuote(SymbolToTrade, marketdata.GetLatestQuoteRequest{})
	if err != nil {
        log.Fatalf("Error getting quote for %s: %v", SymbolToTrade, err)
	}
    quoteJson, err := json.Marshal(quote)
    if err == nil {
        log.Printf("Quote for %s: %v\n", SymbolToTrade, string(quoteJson))
    }

	currentPrice := quote.AskPrice
    if quote.AskPrice == 0 {
        currentPrice = quote.BidPrice
    }
	takeProfitPrice := decimal.NewFromFloat((currentPrice * 0.97))
	stopLossPrice := decimal.NewFromFloat((currentPrice * 1.02))
    log.Printf("Current Price: %v, Take Profit: %v, Stop Loss %v", currentPrice, takeProfitPrice, stopLossPrice)

	orderReq := alpaca.PlaceOrderRequest{
		Symbol:      SymbolToTrade,
		Notional:    &notional,
		Side:        alpaca.Sell,   // Order side: Buy or Sell
		Type:        alpaca.Market, // Order type: Market or Limit
		TimeInForce: alpaca.Day,    // Time in force: Day, GTC, etc.
		OrderClass:  alpaca.Bracket,
		TakeProfit: &alpaca.TakeProfit{
			LimitPrice: &takeProfitPrice,
		},
		StopLoss: &alpaca.StopLoss{
			StopPrice: &stopLossPrice,
		},
        PositionIntent: alpaca.SellToOpen,
	}

	// Submit the order
	order, err := client.Client.PlaceOrder(orderReq)
	if err != nil {
		log.Println("Error placing order:", err)
		return
	}
	orderJSON, err := json.Marshal(order)
	if err == nil {
		log.Printf("Order placed for %s: %+v", SymbolToTrade, string(orderJSON))
	}
}
