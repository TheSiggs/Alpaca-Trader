package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AlpacaConfig  AlpacaConfig
	PolygonConfig PolygonConfig
	Env           string
}

type AlpacaConfig struct {
	APIKey        string
	APISecret     string
	BaseURL       string
	MarketBaseURL string
}

const AlpacaBaseURLSandbox = "https://paper-api.alpaca.markets"
const AlpacaBaseURLProduction = "https://api.alpaca.markets"

const AlpacaMarketBaseURLProduction = "https://data.alpaca.markets"

type PolygonConfig struct {
	APIKey string
}

func Setup() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	env := os.Getenv("ENV")
	switch env {
	case "production":
		env = "production"
        break;
    case "developement":
        env = "developement"
        break;
	case "test":
		env = "test"
        break;
    default:
        env = "developement"
        break;
	}
	log.Printf("Runnng in %s mode\n", env)

	var baseURL string
	if env == "production" {
		baseURL = AlpacaBaseURLProduction
	} else {
		baseURL = AlpacaBaseURLSandbox
	}

	alpacaAPIKey := os.Getenv("ALPACA_API_KEY")
	if alpacaAPIKey == "" {
		log.Fatalf("Missing environment variable: ALPACA_API_KEY")
	}

	alpacaAPISecret := os.Getenv("ALPACA_API_SECRET")
	if alpacaAPISecret == "" {
		log.Fatalf("Missing environment variable: ALPACA_API_SECRET")
	}

	polygonAPIKey := os.Getenv("POLYGON_API_KEY")
	if polygonAPIKey == "" {
		log.Fatalf("Missing environment variable: POLYGON_API_KEY")
	}

	AlpacaConfig := AlpacaConfig{
		BaseURL:       baseURL,
		APIKey:        alpacaAPIKey,
		APISecret:     alpacaAPISecret,
		MarketBaseURL: AlpacaMarketBaseURLProduction,
	}

	PolygonConfig := PolygonConfig{
		APIKey: polygonAPIKey,
	}

	Config := Config{
		AlpacaConfig:  AlpacaConfig,
		PolygonConfig: PolygonConfig,
		Env:           env,
	}
	return Config
}
