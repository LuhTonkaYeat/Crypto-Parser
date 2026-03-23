package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CryptoPrice struct {
	Coin        string             `json:"coin"`
	Name        string             `json:"name"`
	Prices      map[string]float64 `json:"prices"`
	Changes     map[string]float64 `json:"changes"`
	LastUpdated string             `json:"last_updated"`
}

func GetPricesForCoins(coins []string, currencies []string) (map[string]CryptoPrice, error) {
	coinIDs := strings.Join(coins, ",")
	vsCurrencies := strings.Join(currencies, ",")

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s&include_24hr_change=true", coinIDs, vsCurrencies)

	client := &http.Client{Timeout: 10 * time.Second}
	response, error := client.Get(url)
	if error != nil {
		return nil, fmt.Errorf("Network error: %w", error)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("CoinGecko returned status: %d", response.StatusCode)
	}

	var raw map[string]map[string]interface{}
	error = json.NewDecoder(response.Body).Decode(&raw)
	if error != nil {
		return nil, fmt.Errorf("JSON parse error: %w", error)
	}

	result := make(map[string]CryptoPrice)

	for coinID, data := range raw {
		prices := make(map[string]float64)
		changes := make(map[string]float64)

		for key, value := range data {
			switch v := value.(type) {
			case float64:
				if strings.HasSuffix(key, "_24h_change") {
					changes[key] = v
				} else {
					prices[key] = v
				}
			}
		}

		var name string
		switch coinID {
		case "bitcoin":
			name = "Bitcoin"
		case "the-open-network":
			name = "TON"
		case "solana":
			name = "Solana"
		default:
			name = coinID
		}

		result[coinID] = CryptoPrice{
			Coin:        coinID,
			Name:        name,
			Prices:      prices,
			Changes:     changes,
			LastUpdated: time.Now().Format(time.RFC3339),
		}
	}

	return result, nil
}

func GetMainCryptoPrices() (map[string]CryptoPrice, error) {
	coins := []string{
		"bitcoin",
		"the-open-network",
		"ethereum",
		"solana",
	}

	currencies := []string{
		"usd",
		"eur",
		"rub",
	}

	return GetPricesForCoins(coins, currencies)
}
