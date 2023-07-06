package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CurrencyRate struct {
	StartDate     string  `json:"StartDate"`
	TimeSign      string  `json:"TimeSign"`
	CurrencyCode  string  `json:"CurrencyCode"`
	CurrencyCodeL string  `json:"CurrencyCodeL"`
	Units         int     `json:"Units"`
	Amount        float32 `json:"Amount"`
}

func main() {

	usdRate, err := getUSDRate()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("USD:", usdRate)
}

func getUSDRate() (float32, error) {
	// Create a new HTTP client
	client := http.DefaultClient

	// Send an HTTP GET request
	resp, err := client.Get("https://bank.gov.ua/NBU_Exchange/exchange?json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var decodedRates []CurrencyRate
	err = json.Unmarshal(body, &decodedRates)
	if err != nil {
		return 0, err
	}

	var usdRate CurrencyRate
	for _, currencyRate := range decodedRates {
		if currencyRate.CurrencyCodeL == "USD" {
			usdRate = currencyRate
			break
		}
	}

	return usdRate.Amount, nil
}
