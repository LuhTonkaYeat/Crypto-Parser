package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LuhTonkaYeat/crypto-parser/parser"
)

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
		os.Exit(1)
	}

	userInput := strings.ToLower(os.Args[1])

	coinID, exists := coinMapping[userInput]
	if !exists {
		fmt.Printf("Unknown coin: %s\n", userInput)
		fmt.Println("Available: btc, ton, eth, sol")
		os.Exit(1)
	}

	prices, err := parser.GetPricesForCoins([]string{coinID}, []string{"usd", "eur", "rub"})
	if err != nil {
		log.Fatal("Error:", err)
	}

	if data, ok := prices[coinID]; ok {
		fmt.Printf("\n%s (%s):\n", data.Name, userInput)
		fmt.Printf("USD: $%.2f", data.Prices["usd"])
		if ch, ok := data.Changes["usd_24h_change"]; ok {
			fmt.Printf(" (%.2f%%)", ch)
		}
		fmt.Println()

		fmt.Printf("EUR: €%.2f", data.Prices["eur"])
		if ch, ok := data.Changes["eur_24h_change"]; ok {
			fmt.Printf(" (%.2f%%)", ch)
		}
		fmt.Println()

		fmt.Printf("RUB: ₽%.0f", data.Prices["rub"])
		if ch, ok := data.Changes["rub_24h_change"]; ok {
			fmt.Printf(" (%.2f%%)", ch)
		}
		fmt.Println()

		fmt.Printf("Updated: %s\n", data.LastUpdated[:19])
	}
}
