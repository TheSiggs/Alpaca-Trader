package main

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nettis/alpaca-trader/config"
	"github.com/nettis/alpaca-trader/entities"
	"github.com/nettis/alpaca-trader/strategies"
)

func main() {
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
		BaseURL:   client.Config.AlpacaConfig.MarketBaseURL,
	})

	year, month, day := client.CurrentDate(0)

	client.Client.CloseAllPositions(alpaca.CloseAllPositionsRequest{
		CancelOrders: true,
	})
	switch strategy := client.Config.Stragegy; strategy {
	case "ExdividendShorter":
		strategies.ExdividendShorter(&client, year, month, day)
	case "DividendFlipper":
		strategies.DividendFlipper(&client, year, month, day)
	}
}
