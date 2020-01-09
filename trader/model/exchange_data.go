package model

import (
	"strconv"
	"time"
)

type ExchangeData struct {
	OpenTime     time.Time `json:"-"`
	OpenValue    float64   `json:"open_value"`
	High         float64   `json:"high"`
	Low          float64   `json:"low"`
	CloseValue   float64   `json:"close_value"`
	Volume       float64   `json:"volume"`
	QuoteAverage float64   `json:"quote_asset_volume"`
	Trades       int       `json:"trades"`
	TbBaseAvg    float64   `json:"taker_buy_base_asset_volume"`
	TbQuoteAvg   float64   `json:"taker_buy_quote_asset_volume"`
	Ma5          float64   `json:"ma5"`
	Ma10         float64   `json:"ma10"`
}

func NewExchangeFromSlice(slice []string) *ExchangeData {

	//open_time,open_value,high,low,close_value,volume,quote_asset_volume,trades,taker_buy_base_asset_volume,taker_buy_quote_asset_volume,ma5,ma10,prediction
	//2017-08-17 06:30:00
	openTime, err := time.Parse("2006-01-02 15:04:05", slice[1])
	openValue, err := strconv.ParseFloat(slice[2], 64)
	high, err := strconv.ParseFloat(slice[3], 64)
	low, err := strconv.ParseFloat(slice[4], 64)
	closeValue, err := strconv.ParseFloat(slice[5], 64)
	volume, err := strconv.ParseFloat(slice[6], 64)
	quoteAvg, err := strconv.ParseFloat(slice[7], 64)
	trades, err := strconv.Atoi(slice[8])
	tbBaseAvg, err := strconv.ParseFloat(slice[9], 64)
	tbQuoteAvg, err := strconv.ParseFloat(slice[10], 64)
	ma5, err := strconv.ParseFloat(slice[11], 64)
	ma10, err := strconv.ParseFloat(slice[12], 64)

	if err != nil || openTime.IsZero() {
		panic(err)
	}

	return &ExchangeData{
		OpenTime:     openTime,
		OpenValue:    openValue,
		High:         high,
		Low:          low,
		CloseValue:   closeValue,
		Volume:       volume,
		QuoteAverage: quoteAvg,
		Trades:       trades,
		TbBaseAvg:    tbBaseAvg,
		TbQuoteAvg:   tbQuoteAvg,
		Ma5:          ma5,
		Ma10:         ma10,
	}
}
