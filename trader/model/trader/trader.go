package trader

import (
	"scoing-trader/trader/model/market"
	"scoing-trader/trader/model/market/model"
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
			if decision.EventType != HOLD && decision.Qty > 0 {
				record := TradeRecord{
					Timestamp:   prediction.Timestamp,
					Coin:        coin,
					Event:       decision.EventType,
					Qty:         decision.Qty,
					Value:       prediction.CloseValue,
					Transaction: model.IntToFloat(model.IntFloatMul(decision.Val, decision.Qty*(1+t.Accountant.GetFee()))),
					Profit:      0,
				}

				if decision.EventType == SELL {
					record.Profit += model.IntToFloat(model.IntFloatMul(model.FloatToInt(prediction.CloseValue), decision.Qty*(1-t.Accountant.GetFee())) -
						model.IntFloatMul(decision.Val, decision.Qty*(1+t.Accountant.GetFee())))
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
