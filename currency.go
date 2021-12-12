package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

var currencyCache = make(map[string]float64)
var currencyCacheExpireTime = time.Now().UnixMilli()

type Currency struct {
	Eur float64 `json:"EUR"`
	Usd float64 `json:"USD"`
	Gbp float64 `json:"GBP"`
}

func getCurrencyData(provider string, ch chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(provider)
	if err != nil {
		log.Println("Error getting remote data")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response data")
	}
	ch <- b
}

func updateCurrencyCache(providers []string) {
	ch := make(chan []byte)
	var provider string
	var wg sync.WaitGroup
	for _, provider = range providers {
		wg.Add(1)
		go getCurrencyData(provider, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for response := range ch {
		var result Currency
		if err := json.Unmarshal(response, &result); err != nil {
			log.Println("Can not unmarshal JSON")
		}
		if currencyCache["USD"] > result.Usd {
			currencyCache["USD"] = result.Usd
		}
		if currencyCache["EUR"] > result.Eur {
			currencyCache["EUR"] = result.Eur
		}
		if currencyCache["GBP"] > result.Gbp {
			currencyCache["GBP"] = result.Gbp
		}
		currencyCacheExpireTime = time.Now().UnixMilli() + (600000)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func getCurrency(currencyType string) (float64, error) {
	if currencyCacheExpireTime < time.Now().UnixMilli() {
		log.Println("Updating cache")
		currencyProviders, err := readLines("./providers")
		if err != nil {
			log.Println("Can't read providers file")
			return 0.0, err
		}
		currencyCache["USD"] = math.MaxFloat64
		currencyCache["EUR"] = math.MaxFloat64
		currencyCache["GBP"] = math.MaxFloat64
		updateCurrencyCache(currencyProviders)
	}
	result, ok := currencyCache[currencyType]
	if !ok || result == math.MaxFloat64 {
		return 0.0, errors.New("not found")
	}
	return result, nil
}
