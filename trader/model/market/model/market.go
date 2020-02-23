package model

type Market interface {
	NewOrder(order OrderRequest) error
	OpenOrders(symbol string) []*OrderResponseFull
	OrderHistory() []*OrderResponseFull
	CancelOrder(orderId string) error
	AccountInformation() AccountInformation
	Balance(asset string) (Balance, error)
	Trades() []*Trade
	UpdateInformation()
	CoinValue(asset string) (float64, error)
	Deposit(asset string, qty float64)
	UpdateCoinValue(asset string, value float64)
}
