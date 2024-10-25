package strategies

import (
	"encoding/json"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nettis/alpaca-trader/entities"
	"github.com/shopspring/decimal"
)

func ExdividendShorter(client *entities.TradingClient, year int, month time.Month, day int) (*alpaca.Order, error) {
	log.Println("Local Time Zone Location:", time.Local)
	acct, err := client.Client.GetAccount()
	if err != nil {
		panic(err)
	}

    SymbolToTrade := client.LargestDividendStock(year, month, day)

	// Buy the stock
	cash := acct.Cash // Amount in dollars you want to invest

	quote, err := client.MarketClient.GetLatestQuote(SymbolToTrade, marketdata.GetLatestQuoteRequest{})
	if err != nil {
		log.Fatalf("Error getting quote for %s: %v", SymbolToTrade, err)
	}
	quoteJson, err := json.Marshal(quote)
	if err == nil {
		log.Printf("Quote for %s: %v\n", SymbolToTrade, string(quoteJson))
	}

	price := quote.AskPrice
	if quote.AskPrice == 0 {
		price = quote.BidPrice
	}

    currentPrice := decimal.NewFromFloat(price)
    takeProfitThreshold := decimal.NewFromFloat(0.98)
	stopLossThreshold := decimal.NewFromFloat(1.02)

	qty := cash.Div(currentPrice).Floor()
	takeProfitPrice := currentPrice.Mul(takeProfitThreshold)
	stopLossPrice := currentPrice.Mul(stopLossThreshold)

	if takeProfitPrice.GreaterThan(currentPrice.Sub(decimal.NewFromFloat(0.01))) {
		takeProfitPrice = currentPrice.Sub(decimal.NewFromInt(1))
	}

	if stopLossPrice.GreaterThan(currentPrice.Add(decimal.NewFromFloat(0.01))) {
		stopLossPrice = currentPrice.Add(decimal.NewFromInt(1))
	}

    takeProfitPrice = takeProfitPrice.Round(2)
    stopLossPrice = stopLossPrice.Round(2)
    currentPrice = currentPrice.Round(2)

    log.Printf("Symbol: %v, Current Price: %v, Take Profit: %v, Stop Loss %v, Qty: %v", SymbolToTrade, currentPrice, takeProfitPrice, stopLossPrice, qty)

	orderReq := alpaca.PlaceOrderRequest{
		Symbol:      SymbolToTrade,
		Qty:         &qty,
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
		return nil, err
	}
	orderJSON, err := json.Marshal(order)
	if err == nil {
		log.Printf("Order placed for %s: %+v", SymbolToTrade, string(orderJSON))
	}
    return order, nil
}
