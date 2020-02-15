package trader

import (
	"scoing-trader/trader/model/market/model"
	"scoing-trader/trader/model/predictor"
)

type Trader struct {
	Config      StrategyConfig
	Accountant  model.Accountant
	Predictor   predictor.Predictor
	Strategy    Strategy
	Records     []TradeRecord
	KeepRecords bool
}

func NewTrader(config StrategyConfig, accountant model.Accountant, predictor predictor.Predictor, strategy Strategy, keepRecords bool) *Trader {
	return &Trader{
		Config:      config,
		Accountant:  accountant,
		Predictor:   predictor,
		Strategy:    strategy,
		Records:     make([]TradeRecord, 0),
		KeepRecords: keepRecords,
	}
}

func (t *Trader) ProcessData(coin string) {
	prediction := t.Predictor.Predict(coin)

	decisionArr := t.Strategy.ComputeDecision(prediction, t.Accountant.GetPositions(coin), t.Accountant.CoinNetWorth(coin),
		t.Accountant.NetWorth(), t.Accountant.GetBalance(), t.Accountant.GetFee())

	if t.KeepRecords {
		for _, decision := range decisionArr {
			if decision.EventType != HOLD && decision.Qty > 0 {
				record := TradeRecord{
					Timestamp:   prediction.Timestamp,
					Coin:        coin,
					Event:       decision.EventType,
					Qty:         decision.Qty,
					Value:       prediction.CloseValue,
					Transaction: decision.Val * decision.Qty * (1 + t.Accountant.GetFee()),
					Profit:      0,
				}

				if decision.EventType == PROFIT_SELL || decision.EventType == LOSS_SELL {
					record.Profit += (prediction.CloseValue * decision.Qty * (1 - t.Accountant.GetFee())) -
						(decision.Val * decision.Qty * (1 + t.Accountant.GetFee()))
				}

				t.Records = append(t.Records, record)
			}
		}
	}

	for _, decision := range decisionArr {
		if decision.EventType == BUY {
			t.Accountant.Buy(coin, decision.Qty)
		} else if decision.EventType == PROFIT_SELL || decision.EventType == LOSS_SELL {
			t.Accountant.Sell(coin, decision.Val, decision.Qty)
		}
	}
}
