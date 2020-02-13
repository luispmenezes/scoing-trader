package model

type OrderType string

const (
	LIMIT             OrderType = "LIMIT"
	MARKET            OrderType = "MARKET"
	STOP_LOSS         OrderType = "STOP_LOSS"
	STOP_LOSS_LIMIT   OrderType = "STOP_LOSS_LIMIT"
	TAKE_PROFIT       OrderType = "TAKE_PROFIT"
	TAKE_PROFIT_LIMIT OrderType = "TAKE_PROFIT_LIMIT"
	LIMIT_MAKER       OrderType = "LIMIT_MAKER"
)

type OrderSide string

const (
	BUY  OrderSide = "BUY"
	SELL OrderSide = "SELL"
)

type OrderTimeInForce string

const (
	GTC OrderTimeInForce = "GTC"
	FOK OrderTimeInForce = "FOK"
	IOC OrderTimeInForce = "IOC"
)

type OrderResponseType string

const (
	ACK    OrderResponseType = "ACK"
	RESULT OrderResponseType = "RESULT"
	FULL   OrderResponseType = "FULL"
)

type OrderResponseStatus string

const (
	NEW       OrderResponseStatus = "NEW"
	FILLED    OrderResponseStatus = "FILLED"
	CANCELLED OrderResponseStatus = "CANCELLED"
)
