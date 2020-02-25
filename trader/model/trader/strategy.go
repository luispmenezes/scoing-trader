package trader

import (
	"fmt"
	"github.com/shopspring/decimal"
	"scoing-trader/trader/model/predictor"
	"time"
)

type Strategy interface {
	ComputeDecision(prediction predictor.Prediction, positions map[string]decimal.Decimal, coinNetWorth decimal.Decimal,
		totalNetWorth decimal.Decimal, coinValue decimal.Decimal, balance decimal.Decimal, fee decimal.Decimal) map[DecisionType]Decision
	BuySize(prediction predictor.Prediction, coinNetWorth decimal.Decimal, totalNetWorth decimal.Decimal, balance decimal.Decimal, fee decimal.Decimal) decimal.Decimal
	SellSize(prediction predictor.Prediction, positionQty decimal.Decimal, coinValue decimal.Decimal) decimal.Decimal
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
	Qty       decimal.Decimal
	Val       decimal.Decimal
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
	Qty         decimal.Decimal
	Value       decimal.Decimal
	Transaction decimal.Decimal
	Profit      decimal.Decimal
}

func (t *TradeRecord) ToString() string {
	strRecord := fmt.Sprintf("## %s : (%s) - %s -- Qty:%s Value:%s$ Trans:%s$", t.Timestamp, t.Event, t.Coin,
		t.Qty, t.Value, t.Transaction)
	if t.Event != BUY {
		strRecord += fmt.Sprintf(" Profit: %s$", t.Profit)
	}
	return strRecord
}
