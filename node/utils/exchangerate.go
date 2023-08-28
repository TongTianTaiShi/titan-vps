package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var rateInfo struct {
	USDRate float32
	ET      time.Time
}

type tianAPIResponse struct {
	Code int32         `json:"code"`
	Data tianAPIResult `json:"result"`
}
type tianAPIResult struct {
	Rate string `json:"money"`
}

func getExchangeRateWithTianAPI() float32 {
	client := &http.Client{}
	// todo
	// resp, err := client.Get("https://api.it120.cc/gooking/forex/rate?fromCode=CNY&toCode=USD")
	resp, err := client.Get("https://apis.tianapi.com/fxrate/index?key=af490d21502b58010c7feef4db2cd14a&fromcoin=USD&tocoin=CNY&money=1")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var ExchangeRateRsp tianAPIResponse
	err = json.Unmarshal(body, &ExchangeRateRsp)
	if err != nil {
		return 0
	}
	distFloat, _ := strconv.ParseFloat(ExchangeRateRsp.Data.Rate, 32)
	return float32(distFloat)
}

type exchangeRateResponse struct {
	Rates struct {
		CNY float32 `json:"CNY"`
	} `json:"rates"`
}

func getExchangeRateWithExchangeRateAPI() float32 {
	client := &http.Client{}
	// todo
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

// GetUSDRate get rate
func GetUSDRate() float32 {
	if rateInfo.USDRate == 0 || time.Now().After(rateInfo.ET) {
		rateInfo.USDRate = getExchangeRateWithExchangeRateAPI()
		rateInfo.ET = time.Now().Add(time.Hour * 2)
	}

	if rateInfo.USDRate == 0 {
		rateInfo.USDRate = 7.2673
	}

	return rateInfo.USDRate
}
