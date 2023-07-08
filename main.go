package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
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
	apiSecret := os.Getenv("GA_API_CLIENT_ID")
	measurementId := os.Getenv("GA_MEASUREMENT_ID")

	for {
		usdRate, err := getUSDRate()

		if err != nil {
			fmt.Println("Error:", err)
		}

		sendToGA(apiSecret, measurementId, usdRate)

		log.Println("USD rate", usdRate, "was sent to GA")
		time.Sleep(60 * time.Minute)
	}
}

func getUSDRate() (float32, error) {
	client := http.DefaultClient

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

func sendToGA(apiSecret string, measurementId string, usdRate float32) {
	gaURL := buildGAURL(apiSecret, measurementId)

	body, err := buildRequestBody(usdRate)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	request, err := http.NewRequest("POST", gaURL, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
}

func buildGAURL(apiSecret string, measurementId string) string {
	baseURL := "https://www.google-analytics.com"
	resource := "/mp/collect"

	params := url.Values{}
	params.Add("api_secret", apiSecret)
	params.Add("measurement_id", measurementId)

	u, _ := url.ParseRequestURI(baseURL)

	u.Path = resource
	u.RawQuery = params.Encode()

	return u.String()
}

func buildRequestBody(usdRate float32) ([]byte, error) {
	data := map[string]interface{}{
		"client_id":            "689844021.1688749428",
		"non_personalized_ads": true,
		"events": []interface{}{
			map[string]interface{}{
				"name": "usd_rate_checked",
				"params": map[string]interface{}{
					"value": usdRate,
				},
			},
		},
	}

	return json.Marshal(data)
}
