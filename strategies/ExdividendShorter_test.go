package strategies

import (
	"testing"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nettis/alpaca-trader/config"
	"github.com/nettis/alpaca-trader/entities"
)

func TestOrders(t *testing.T) {
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

    for days := 0; days > -10; days-- {
        currentDate := time.Now()
        currentDate = currentDate.AddDate(0, 0, days)
        year, month, day := currentDate.Date()
        order, err := ExdividendShorter(&client, year, month, day)
        if err != nil {
            t.Fatal(err)
        }

        err = client.Client.CancelOrder(order.ID)
        if err != nil {
            t.Fatal(err)
        }
    }
}
