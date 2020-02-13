package trader

import (
	"scoing-trader/trader/model/market"
	"scoing-trader/trader/model/predictor"
)

type Trader struct {
	Config      StrategyConfig
	Wallet      market.Wallet
	Predictor   predictor.Predictor
	Strategy    Strategy
	Records     []TradeRecord
	KeepRecords bool
}

func NewTrader(config StrategyConfig, wallet market.Wallet, predictor predictor.Predictor, strategy Strategy, keepRecords bool) *Trader {
	return &Trader{
		Config:      config,
		Wallet:      wallet,
		Predictor:   predictor,
		Strategy:    strategy,
		Records:     make([]TradeRecord, 0),
		KeepRecords: keepRecords,
	}
}

func (t *Trader) ProcessData(coin string) {
	prediction := t.Predictor.Predict(coin)

	decisionArr := t.Strategy.ComputeDecision(prediction, t.Wallet.GetPositions(coin), t.Wallet.CoinNetWorth(coin),
		t.Wallet.NetWorth(), t.Wallet.GetBalance(), t.Wallet.GetFee())

	if t.KeepRecords {
		for _, decision := range decisionArr {
			if decision.EventType != HOLD && decision.Qty > 0 {
				record := TradeRecord{
					Timestamp:   prediction.Timestamp,
					Coin:        coin,
					Event:       decision.EventType,
					Qty:         decision.Qty,
					Value:       prediction.CloseValue,
					Transaction: decision.Val * decision.Qty * (1 + t.Wallet.GetFee()),
					Profit:      0,
				}

				if decision.EventType == PROFIT_SELL || decision.EventType == LOSS_SELL {
					record.Profit += (prediction.CloseValue * decision.Qty * (1 - t.Wallet.GetFee())) -
						(decision.Val * decision.Qty * (1 + t.Wallet.GetFee()))
				}

				t.Records = append(t.Records, record)
			}
		}
	}

	for _, decision := range decisionArr {
		if decision.EventType == BUY {
			t.Wallet.Buy(coin, decision.Qty)
		} else if decision.EventType == PROFIT_SELL || decision.EventType == LOSS_SELL {
			t.Wallet.Sell(coin, decision.Val, decision.Qty)
		}
	}
}
