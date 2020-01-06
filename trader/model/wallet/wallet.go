package wallet

import "time"

type Wallet interface {
	Buy(coin string, quantity float64)
	Sell(coin string, buyValue float64, quantity float64)
	UpdateCoinValue(coin string, value float64, timestamp time.Time)
	GetBalance() float64
	GetPositions(coin string) map[float64]float64
	TotalPositionValue() float64
	NetWorth() float64
	GetDailyNetWorth() []float64
}