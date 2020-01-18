package trader

import (
	"super-trader/trader/model/predictor"
)

type Strategy interface {
	ComputeDecision(prediction predictor.Prediction, positions map[float64]float64, coinNetWorth float64,
		totalNetWorth float64, balance float64, fee float64) Decision
	BuySize(prediction predictor.Prediction, coinNetWorth float64, totalNetWorth float64, balance float64, fee float64) float64
	SellSize(prediction predictor.Prediction, positionQty float64) float64
}

type Decision struct {
	EventType DecisionType
	Qty       float64
	Val       float64
}

type DecisionType string

const (
	BUY  DecisionType = "BUY"
	SELL DecisionType = "SELL"
	HOLD DecisionType = "HOLD"
)
