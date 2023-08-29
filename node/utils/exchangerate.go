package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var exchangeRateInfo struct {
	USDRate float32
	ET      time.Time
}

// tianAPIResponse represents the response structure from the Tian API.
type tianAPIResponse struct {
	Code int32         `json:"code"`
	Data tianAPIResult `json:"result"`
}

// tianAPIResult represents the result structure in the Tian API response.
type tianAPIResult struct {
	Rate string `json:"money"`
}

// getExchangeRateWithTianAPI fetches the exchange rate from the Tian API.
func getExchangeRateWithTianAPI() float32 {
	client := &http.Client{}
	// TODO: Replace the URL with the correct one.
	resp, err := client.Get("https://apis.tianapi.com/fxrate/index?key=af490d21502b58010c7feef4db2cd14a&fromcoin=USD&tocoin=CNY&money=1")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var exchangeRateRsp tianAPIResponse
	err = json.Unmarshal(body, &exchangeRateRsp)
	if err != nil {
		return 0
	}
	distFloat, _ := strconv.ParseFloat(exchangeRateRsp.Data.Rate, 32)
	return float32(distFloat)
}

// exchangeRateResponse represents the response structure from the ExchangeRate API.
type exchangeRateResponse struct {
	Rates struct {
		CNY float32 `json:"CNY"`
	} `json:"rates"`
}

// getExchangeRateWithExchangeRateAPI fetches the exchange rate from the ExchangeRate API.
func getExchangeRateWithExchangeRateAPI() float32 {
	client := &http.Client{}
	// TODO: Replace the URL with the correct one.
	resp, err := client.Get("https://api.exchangerate-api.com/v4/latest/USD")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var r exchangeRateResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return r.Rates.CNY
}

// GetUSDRate returns the USD to CNY exchange rate.
func GetUSDRate() float32 {
	if exchangeRateInfo.USDRate == 0 || time.Now().After(exchangeRateInfo.ET) {
		exchangeRateInfo.USDRate = getExchangeRateWithExchangeRateAPI()
		exchangeRateInfo.ET = time.Now().Add(time.Hour * 2)
	}

	if exchangeRateInfo.USDRate == 0 {
		exchangeRateInfo.USDRate = 7.2673
	}

	return exchangeRateInfo.USDRate
}
