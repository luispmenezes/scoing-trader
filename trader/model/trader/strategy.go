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
	BuyConf   float64
	SellConf  float64
	DebugText string
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
	BuyConf     float64
	SellConf    float64
	DebugText   string
}

func (t *TradeRecord) ToString() string {
	qty, _ := t.Qty.Float64()
	value, _ := t.Value.Float64()
	transaction, _ := t.Transaction.Float64()

	strRecord := fmt.Sprintf("## %s : (%s) - %s -- Qty:%.4f Value:%.4f$ Trans:%.4f$", t.Timestamp, t.Event, t.Coin,
		qty, value, transaction)
	if t.Event == SELL {
		profit, _ := t.Profit.Float64()
		strRecord += fmt.Sprintf(" Profit: %.4f$", profit)
	} else if t.Event == HOLD {
		strRecord = fmt.Sprintf("## HOLD: Coin: %s BuyConfidence: %.2f/2 SellConfidence: %.2f/-2", t.Coin, t.BuyConf, t.SellConf)
		strRecord += t.DebugText
	}

	return strRecord
}
