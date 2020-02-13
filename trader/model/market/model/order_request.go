package model

type OrderRequest struct {
	Symbol        string
	Side          OrderSide
	Type          OrderType
	Timestamp     int64
	TimeInForce   OrderTimeInForce
	Quantity      float64
	QuoteOrderQty float64
	Price         float64
	ClientOrderId string
	StopPrice     float64
	IcebergQty    float64
	ResponseType  OrderResponseType
}
