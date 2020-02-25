package trader

import (
	"github.com/shopspring/decimal"
	"scoing-trader/trader/model/market"
	"scoing-trader/trader/model/predictor"
)

type Trader struct {
	Accountant  market.Accountant
	Predictor   predictor.Predictor
	Strategy    Strategy
	Records     []TradeRecord
	KeepRecords bool
}

func NewTrader(accountant market.Accountant, predictor predictor.Predictor, strategy Strategy, keepRecords bool) *Trader {
	return &Trader{
		Accountant:  accountant,
		Predictor:   predictor,
		Strategy:    strategy,
		Records:     make([]TradeRecord, 0),
		KeepRecords: keepRecords,
	}
}

func (t *Trader) ProcessData(coin string) {
	prediction := t.Predictor.Predict(coin)

	decisionArr := t.Strategy.ComputeDecision(prediction, t.Accountant.GetPositions(coin), t.Accountant.AssetValue(coin),
		t.Accountant.NetWorth(), t.Accountant.AssetValues[coin], t.Accountant.GetBalance(), t.Accountant.GetFee())

	if t.KeepRecords {
		for _, decision := range decisionArr {
			if decision.EventType != HOLD && decision.Qty.LessThan(decimal.Zero) {
				record := TradeRecord{
					Timestamp:   prediction.Timestamp,
					Coin:        coin,
					Event:       decision.EventType,
					Qty:         decision.Qty,
					Value:       decimal.NewFromFloat(prediction.CloseValue),
					Transaction: decision.Val.Mul(decision.Qty.Mul(decimal.NewFromInt(1).Add(t.Accountant.Fee))),
					Profit:      decimal.Zero,
				}

				if decision.EventType == SELL {
					//TODO: NEED TO ACCOUNT FOR MULTIPLE POSITIONS PER SALE
					record.Profit = record.Profit.Add(decimal.NewFromFloat(prediction.CloseValue).
						Mul(decision.Qty.Mul(decimal.NewFromInt(1).Sub(t.Accountant.Fee))).
						Sub(decision.Val.Mul(decision.Qty).Mul(decimal.NewFromInt(1).Add(t.Accountant.Fee))))
				}

				t.Records = append(t.Records, record)
			}
		}
	}

	for _, decision := range decisionArr {
		if decision.EventType == BUY {
			err := t.Accountant.Buy(coin, decision.Qty)
			if err != nil {
				panic(err)
			}
		} else if decision.EventType == SELL {
			err := t.Accountant.Sell(coin, decision.Qty)
			if err != nil {
				panic(err)
			}
		}
	}
}
