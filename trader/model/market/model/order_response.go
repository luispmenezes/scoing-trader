package model

import "github.com/shopspring/decimal"

type OrderResponseAck struct {
	Symbol          string
	OrderId         float64
	OrderListId     float64
	ClientOrderId   string
	TransactionTime float64
}

type OrderResponseResult struct {
	Symbol              string
	OrderId             float64
	OrderListId         float64
	ClientOrderId       string
	TransactionTime     float64
	Price               float64
	OrigQty             float64
	ExecutedQty         float64
	CummulativeQuoteQty float64
	Status              OrderResponseStatus
	TimeInForce         OrderTimeInForce
	Type                OrderType
	Side                OrderSide
}

type OrderResponseFull struct {
	Symbol              string
	OrderId             int64
	OrderListId         int64
	ClientOrderId       string
	TransactionTime     int64
	Price               decimal.Decimal
	OrigQty             decimal.Decimal
	ExecutedQty         decimal.Decimal
	CummulativeQuoteQty decimal.Decimal
	Status              OrderResponseStatus
	TimeInForce         OrderTimeInForce
	Type                OrderType
	Side                OrderSide
	Fills               []Fill
}

type Fill struct {
	Price           decimal.Decimal
	Qty             decimal.Decimal
	Commission      decimal.Decimal
	CommissionAsset string
}
