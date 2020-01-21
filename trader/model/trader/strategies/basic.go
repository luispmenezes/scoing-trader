package strategies

import (
	"math"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
)

type BasicStrategy struct {
	Config BasicConfig
}

func NewBasicStrategy(slice []float64) *BasicStrategy {
	basicStrategy := &BasicStrategy{Config: BasicConfig{}}
	basicStrategy.Config.FromSlice(slice)
	return basicStrategy
}

func (s *BasicStrategy) ComputeDecision(prediction predictor.Prediction, positions map[float64]float64,
	coinNetWorth float64, totalNetWorth float64, balance float64, fee float64) trader.Decision {

	pred15 := 0.0
	pred60 := 0.0
	pred1440 := 0.0

	if prediction.Pred15 > 0 {
		pred15 = 1
	}
	if prediction.Pred60 > 0 {
		pred60 = 1
	}
	if prediction.Pred1440 > 0 {
		pred1440 = 1
	}

	if ((pred15 * s.Config.BuyPred15Mod) + (pred60 * s.Config.BuyPred60Mod) + (pred1440 * s.Config.BuyPred1440Mod)) > 2 {
		decision := trader.Decision{
			EventType: trader.BUY,
			Coin:      prediction.Coin,
			Qty:       s.BuySize(prediction, coinNetWorth, totalNetWorth, balance, fee),
			Val:       prediction.CloseValue,
		}
		if decision.Qty > 0 {
			return decision
		}
	}

	if ((pred15 * s.Config.SellPred15Mod) + (pred60 * s.Config.SellPred60Mod) + (pred1440 * s.Config.SellPred1440Mod)) > 2 {
		for val, qty := range positions {
			currentProfit := 1 - (val / prediction.CloseValue)
			if currentProfit < s.Config.StopLoss {
				decision := trader.Decision{
					EventType: trader.LOSS_SELL,
					Coin:      prediction.Coin,
					Qty:       s.SellSize(prediction, qty),
					Val:       val,
				}
				if decision.Qty > 0 {
					return decision
				}
			}
		}
	}

	for val, qty := range positions {
		currentProfit := 1 - (val / prediction.CloseValue)
		if currentProfit > s.Config.ProfitCap {
			decision := trader.Decision{
				EventType: trader.PROFIT_SELL,
				Coin:      prediction.Coin,
				Qty:       s.SellSize(prediction, qty),
				Val:       val,
			}
			if decision.Qty > 0 {
				return decision
			}
		}
	}

	return trader.Decision{
		EventType: trader.HOLD,
		Coin:      prediction.Coin,
		Qty:       0,
		Val:       0,
	}
}

func (s *BasicStrategy) BuySize(prediction predictor.Prediction, coinNetWorth float64, totalNetWorth float64,
	balance float64, fee float64) float64 {

	maxCoinNetWorth := totalNetWorth * 0.3
	maxTransaction := totalNetWorth * 0.05

	if maxCoinNetWorth-coinNetWorth >= 10 && balance >= 10*(1+fee) {
		transaction := math.Max(10, math.Min(maxTransaction, (maxCoinNetWorth-coinNetWorth)*s.Config.BuyQtyMod))
		transactionWFee := transaction * (1 + fee)

		if transactionWFee < balance {
			return transaction / prediction.CloseValue
		} else {
			return balance / (prediction.CloseValue * (1 + fee))
		}
	} else {
		return 0.0
	}
}

func (s *BasicStrategy) SellSize(prediction predictor.Prediction, positionQty float64) float64 {
	return positionQty * s.Config.SellQtyMod
}
