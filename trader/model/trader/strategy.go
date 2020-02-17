package trader

import (
	"fmt"
	"scoing-trader/trader/model/predictor"
	"time"
)

type Strategy interface {
	ComputeDecision(prediction predictor.Prediction, positions map[int64]float64, coinNetWorth int64,
		totalNetWorth int64, coinValue int64, balance int64, fee float64) map[DecisionType]Decision
	BuySize(prediction predictor.Prediction, coinNetWorth int64, totalNetWorth int64, balance int64, fee float64) float64
	SellSize(prediction predictor.Prediction, positionQty float64, coinValue int64) float64
}

type StrategyConfig interface {
	NumParams() int
	ToSlice() []float64
	FromSlice(slice []float64)
	ParamRanges() ([]float64, []float64)
	RandomFromSlices(a []float64, b []float64)
	RandomizeParam()
}

type Decision struct {
	EventType DecisionType
	Coin      string
	Qty       float64
	Val       int64
}

type DecisionType string

const (
	BUY  DecisionType = "BUY"
	SELL DecisionType = "SELL"
	HOLD DecisionType = "HOLD"
)

type TradeRecord struct {
	Timestamp   time.Time
	Coin        string
	Event       DecisionType
	Qty         float64
	Value       float64
	Transaction float64
	Profit      float64
}

func (t *TradeRecord) ToString() string {
	strRecord := fmt.Sprintf("## %s : (%s) - %s -- Qty:%f Value:%f$ Trans:%f$", t.Timestamp, t.Event, t.Coin,
		t.Qty, t.Value, t.Transaction)
	if t.Event != BUY {
		strRecord += fmt.Sprintf(" Profit: %f$", t.Profit)
	}
	return strRecord
}
