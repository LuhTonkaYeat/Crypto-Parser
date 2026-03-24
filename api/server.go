package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LuhTonkaYeat/crypto-parser/parser"
)

type Cache struct {
	data      map[string]parser.CryptoPrice
	updatedAt time.Time
	mutex     sync.RWMutex
	ttl       time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		data:      make(map[string]parser.CryptoPrice),
		ttl:       ttl,
		updatedAt: time.Time{},
	}
}

func (c *Cache) isStale() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return time.Since(c.updatedAt) > c.ttl
}

func (c *Cache) refresh() error {
	log.Println("Refreshing cache from CoinGecko...")

	prices, err := parser.GetMainCryptoPrices()
	if err != nil {
		return fmt.Errorf("failed to fetch prices: %w", err)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = prices
	c.updatedAt = time.Now()

	log.Printf("✅ Cache refreshed (%d coins)", len(c.data))
	return nil
}

func (c *Cache) get() (map[string]parser.CryptoPrice, error) {
	if c.isStale() {
		if err := c.refresh(); err != nil {
			c.mutex.RLock()
			defer c.mutex.RUnlock()
			if len(c.data) > 0 {
				log.Println("Using cache data (irrelevant)")
				return c.data, nil
			}
			return nil, err
		}
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data, nil
}

func main() {
	cache := NewCache(1 * time.Minute)

	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		prices, err := cache.get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})

	http.HandleFunc("/price/", func(w http.ResponseWriter, r *http.Request) {
		coin := r.URL.Path[len("/price/"):]
		if coin == "" {
			http.Error(w, "coin required", http.StatusBadRequest)
			return
		}

		prices, err := cache.get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for id, data := range prices {
			if id == coin || data.Name == coin {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(data)
				return
			}
		}
		http.Error(w, fmt.Sprintf("coin '%s' not found", coin), http.StatusNotFound)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "ok",
			"cache_ttl":  cache.ttl.String(),
			"cached_at":  cache.updatedAt.Format(time.RFC3339),
			"cached_cnt": len(cache.data),
		})
	})

	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Printf("GET /prices         - all crypto prices")
	log.Printf("GET /price/bitcoin  - price for specific coin")
	log.Printf("GET /health         - server health")
	log.Fatal(http.ListenAndServe(port, nil))
}
