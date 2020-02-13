package market

import "scoing-trader/trader/model/market/model"

type Market interface {
	NewOrder(order model.OrderRequest)
	OpenOrders() []model.OrderResponseFull
	OrderHistory() []model.OrderResponseFull
	CancelOrder(orderId string)
	AccountInformation() model.AccountInformation
	Trades() []model.Trade
	ToString() string
}
