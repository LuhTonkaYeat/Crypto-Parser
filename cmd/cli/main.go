package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type CryptoPrice struct {
	Coin        string             `json:"coin"`
	Name        string             `json:"name"`
	Prices      map[string]float64 `json:"prices"`
	Changes     map[string]float64 `json:"changes"`
	LastUpdated string             `json:"last_updated"`
}

var coinMapping = map[string]string{
	"btc":      "bitcoin",
	"bitcoin":  "bitcoin",
	"ton":      "the-open-network",
	"eth":      "ethereum",
	"ethereum": "ethereum",
	"sol":      "solana",
	"solana":   "solana",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: crypto-cli <coin>")
		fmt.Println("Available: btc, bitcoin, ton, eth, ethereum, sol, solana")
		fmt.Println("Example: crypto-cli btc")
		fmt.Println("\nMake sure the server is running:")
		fmt.Println("   go run api/server.go")
		os.Exit(1)
	}

	userInput := strings.ToLower(os.Args[1])

	coinName, exists := coinMapping[userInput]
	if !exists {
		fmt.Printf("Unknown coin: %s\n", userInput)
		fmt.Println("Available: btc, ton, eth, sol")
		os.Exit(1)
	}

	apiURL := fmt.Sprintf("http://localhost:8080/price/%s", coinName)
	fmt.Printf("Requesting: %s\n", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("\nError connecting to API: %v\n", err)
		fmt.Println("\nMake sure the server is running:")
		fmt.Println("   cd ~/crypto-parser")
		fmt.Println("   go run api/server.go")
		os.Exit(1)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		fmt.Printf("Coin '%s' not found in API\n", coinName)
		os.Exit(1)
	default:
		fmt.Printf("API error: %s (status %d)\n", resp.Status, resp.StatusCode)
		os.Exit(1)
	}

	var data CryptoPrice
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("%s (%s)\n", data.Name, userInput)
	fmt.Println(strings.Repeat("-", 50))

	if usd, ok := data.Prices["usd"]; ok {
		fmt.Printf("💶 USD: $%.2f", usd)
		if ch, ok := data.Changes["usd_24h_change"]; ok {
			emoji := "📈"
			if ch < 0 {
				emoji = "📉"
			}
			fmt.Printf("  %s %.2f%%", emoji, ch)
		}
		fmt.Println()
	}

	if eur, ok := data.Prices["eur"]; ok {
		fmt.Printf("💶 EUR: €%.2f", eur)
		if ch, ok := data.Changes["eur_24h_change"]; ok {
			emoji := "📈"
			if ch < 0 {
				emoji = "📉"
			}
			fmt.Printf("  %s %.2f%%", emoji, ch)
		}
		fmt.Println()
	}

	if rub, ok := data.Prices["rub"]; ok {
		fmt.Printf("💷 RUB: ₽%.0f", rub)
		if ch, ok := data.Changes["rub_24h_change"]; ok {
			emoji := "📈"
			if ch < 0 {
				emoji = "📉"
			}
			fmt.Printf("  %s %.2f%%", emoji, ch)
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Last updated: %s\n", data.LastUpdated[:19])
	fmt.Println(strings.Repeat("=", 50))
}
