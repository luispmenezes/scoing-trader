package trader

import (
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/wallet"
)

type Trader struct {
	Config      TraderConfig
	Wallet      wallet.Wallet
	Predictor   predictor.Predictor
	Strategy    Strategy
	Decisions   []Decision
	KeepRecords bool
}

func NewTrader(config TraderConfig, wallet wallet.Wallet, predictor predictor.Predictor, strategy Strategy, keepRecords bool) *Trader {
	return &Trader{
		Config:      config,
		Wallet:      wallet,
		Predictor:   predictor,
		Strategy:    strategy,
		Decisions:   make([]Decision, 0),
		KeepRecords: keepRecords,
	}
}

func (t *Trader) ProcessData(coin string) {
	prediction := t.Predictor.Predict(coin)

	decision := t.Strategy.ComputeDecision(prediction, t.Wallet.GetPositions(coin), t.Wallet.CoinNetWorth(coin),
		t.Wallet.NetWorth(), t.Wallet.GetBalance(), t.Wallet.GetFee())

	switch decision.EventType {
	case BUY:
		t.Wallet.Buy(coin, decision.Qty)
		break
	case SELL:
		t.Wallet.Sell(coin, decision.Val, decision.Qty)
		break
	}
	if t.KeepRecords {
		if decision.EventType != HOLD && decision.Qty > 0 {
			t.Decisions = append(t.Decisions, decision)
		}
	}
}
