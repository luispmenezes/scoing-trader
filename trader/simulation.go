package trader

import (
	"super-trader/trader/model"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
	"super-trader/trader/model/wallet"
)

type Simulation struct {
	ExchangeData map[string][]model.ExchangeData
	Trader       trader.Trader
}

func NewSimulation(exchangeData map[string][]model.ExchangeData, predictions map[string][]predictor.Prediction,
	config trader.TraderConfig, initialBalance float64, fee float64, uncertainty float64, logging bool) *Simulation {
	return &Simulation{
		ExchangeData: exchangeData,
		Trader: *trader.NewTrader(config,
			wallet.NewSimulatedWallet(initialBalance, fee),
			predictor.NewSimulatedPredictor(predictions, uncertainty), logging),
	}
}

func (s *Simulation) Run() {
	for coin, data := range s.ExchangeData {
		for _, dataEntry := range data {
			s.Trader.ProcessData(dataEntry, coin)
		}
	}
}
