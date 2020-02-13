package market

import "time"

type Wallet interface {
	Buy(coin string, quantity float64)
	Sell(coin string, buyValue float64, quantity float64)
	UpdateCoinValue(coin string, value float64, timestamp time.Time)
	GetBalance() float64
	GetPositions(coin string) map[float64]float64
	GetFee() float64
	TotalPositionValue() float64
	NetWorth() float64
	CoinNetWorth(coin string) float64
	ToString() string
}
