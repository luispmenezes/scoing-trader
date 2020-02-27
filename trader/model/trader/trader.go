package trader

import (
	"github.com/shopspring/decimal"
	"scoing-trader/trader/model/market"
	"scoing-trader/trader/model/predictor"
)

type Trader struct {
	Accountant       market.Accountant
	Predictor        predictor.Predictor
	Strategy         Strategy
	Records          []TradeRecord
	KeepRecords      bool
	OnlyTransactions bool
}

func NewTrader(accountant market.Accountant, predictor predictor.Predictor, strategy Strategy, keepRecords bool, onlyTransactions bool) *Trader {
	return &Trader{
		Accountant:       accountant,
		Predictor:        predictor,
		Strategy:         strategy,
		Records:          make([]TradeRecord, 0),
		KeepRecords:      keepRecords,
		OnlyTransactions: onlyTransactions,
	}
}

func (t *Trader) ProcessData(coin string) {
	prediction := t.Predictor.Predict(coin)

	decisionArr := t.Strategy.ComputeDecision(prediction, t.Accountant.GetPositions(coin), t.Accountant.AssetValue(coin),
		t.Accountant.NetWorth(), t.Accountant.AssetValues[coin], t.Accountant.GetBalance(), t.Accountant.GetFee())

	for _, decision := range decisionArr {
		var transaction decimal.Decimal
		var profit decimal.Decimal
		var err error

		if decision.EventType == BUY {
			transaction, err = t.Accountant.Buy(coin, decision.Qty)
			if err != nil {
				panic(err)
			}
		} else if decision.EventType == SELL {
			transaction, profit, err = t.Accountant.Sell(coin, decision.Qty)
			if err != nil {
				panic(err)
			}
		}

		if t.KeepRecords {
			if !t.OnlyTransactions || decision.EventType != HOLD {
				record := TradeRecord{
					Timestamp:   prediction.Timestamp,
					Coin:        coin,
					Event:       decision.EventType,
					Qty:         decision.Qty,
					Value:       decimal.NewFromFloat(prediction.CloseValue),
					Transaction: transaction,
					Profit:      profit,
					BuyConf:     decision.BuyConf,
					SellConf:    decision.SellConf,
					DebugText:   decision.DebugText,
				}

				t.Records = append(t.Records, record)
			}
		}
	}
}
