package model

import "github.com/shopspring/decimal"

type Market interface {
	NewOrder(order OrderRequest) error
	OpenOrders(symbol string) []*OrderResponseFull
	OrderHistory() []*OrderResponseFull
	CancelOrder(orderId string) error
	AccountInformation() AccountInformation
	Balance(asset string) (Balance, error)
	Trades() []*Trade
	UpdateInformation()
	CoinValue(asset string) (decimal.Decimal, error)
	Deposit(asset string, qty decimal.Decimal)
	UpdateCoinValue(asset string, value decimal.Decimal)
}
