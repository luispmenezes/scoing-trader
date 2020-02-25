package model

import "github.com/shopspring/decimal"

type OrderRequest struct {
	Symbol        string
	Side          OrderSide
	Type          OrderType
	Timestamp     int64
	TimeInForce   OrderTimeInForce
	Quantity      decimal.Decimal
	QuoteOrderQty decimal.Decimal
	Price         decimal.Decimal
	ClientOrderId string
	StopPrice     decimal.Decimal
	IcebergQty    decimal.Decimal
	ResponseType  OrderResponseType
}
