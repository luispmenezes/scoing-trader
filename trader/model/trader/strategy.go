package trader

import (
	"fmt"
	"super-trader/trader/model/predictor"
	"time"
)

type Strategy interface {
	ComputeDecision(prediction predictor.Prediction, positions map[float64]float64, coinNetWorth float64,
		totalNetWorth float64, balance float64, fee float64) Decision
	BuySize(prediction predictor.Prediction, coinNetWorth float64, totalNetWorth float64, balance float64, fee float64) float64
	SellSize(prediction predictor.Prediction, positionQty float64) float64
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
	Val       float64
}

type DecisionType string

const (
	BUY         DecisionType = "BUY"
	PROFIT_SELL DecisionType = "PROFIT_SELL"
	LOSS_SELL   DecisionType = "LOSS_SELL"
	HOLD        DecisionType = "HOLD"
)

type TradeRecord struct {
	Timestamp   time.Time
	Coin        string
	Event       DecisionType
	Qty         float64
	Value       float64
	Transaction float64
	Profit      float64
	NetWorth    float64
	Balance     float64
}

func (t *TradeRecord) ToString() string {
	return fmt.Sprintf("%s : %s (%s) Qty:%f Value%f Trans:%f Profit: %f NW:%f Balance:%f", t.Timestamp, t.Coin,
		t.Event, t.Qty, t.Value, t.Transaction, t.Profit, t.NetWorth, t.Balance)
}
