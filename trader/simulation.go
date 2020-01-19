package trader

import (
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"super-trader/trader/model/trader/strategies"
	"super-trader/trader/model/wallet"
)

type Simulation struct {
	Predictions []predictor.Prediction
	Trader      trader.Trader
	Logging     bool
}

func NewSimulation(predictions []predictor.Prediction, config trader.StrategyConfig, initialBalance float64, fee float64,
	uncertainty float64, keepRecords bool) *Simulation {
	return &Simulation{
		Predictions: predictions,
		Trader: *trader.NewTrader(config,
			wallet.NewSimulatedWallet(initialBalance, fee),
			predictor.NewSimulatedPredictor(uncertainty),
			strategies.NewBasicStrategy(config.ToSlice()), keepRecords),
		Logging: keepRecords,
	}
}

func (sim *Simulation) Run() {
	for _, pred := range sim.Predictions {
		sim.Trader.Wallet.UpdateCoinValue(pred.Coin, pred.CloseValue, pred.Timestamp)
		sim.Trader.Predictor.SetNextPrediction(pred)
		sim.Trader.ProcessData(pred.Coin)
	}
}
