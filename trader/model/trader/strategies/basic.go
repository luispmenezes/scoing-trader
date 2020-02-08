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
	coinNetWorth float64, totalNetWorth float64, balance float64, fee float64) []trader.Decision {

	var decisionArr []trader.Decision

	pred15 := 0.0
	pred60 := 0.0
	pred1440 := 0.0

	if prediction.Pred5 > 0.01 {
		pred15 = 1
	} else if prediction.Pred5 < -0.01 {
		pred15 = -1
	}

	if prediction.Pred10 > 0.01 {
		pred60 = 1
	} else if prediction.Pred10 < -0.01 {
		pred60 = -1
	}

	if prediction.Pred100 > 0.01 {
		pred1440 = 1
	} else if prediction.Pred100 < -0.01 {
		pred1440 = -1
	}

	if ((pred15 * s.Config.BuyPred5Mod) + (pred60 * s.Config.BuyPred10Mod) + (pred1440 * s.Config.BuyPred100Mod)) > 2 {
		decision := trader.Decision{
			EventType: trader.BUY,
			Coin:      prediction.Coin,
			Qty:       s.BuySize(prediction, coinNetWorth, totalNetWorth, balance, fee),
			Val:       prediction.CloseValue,
		}
		if decision.Qty > 0 {
			decisionArr = append(decisionArr, decision)
		}
	}

	if ((pred15 * s.Config.SellPred5Mod) + (pred60 * s.Config.SellPred10Mod) + (pred1440 * s.Config.SellPred100Mod)) < -2 {
		for val, qty := range positions {
			currentProfit := 1 - ((val / prediction.CloseValue) * (1 - fee))
			if currentProfit < s.Config.StopLoss {
				decision := trader.Decision{
					EventType: trader.LOSS_SELL,
					Coin:      prediction.Coin,
					Qty:       s.SellSize(prediction, qty),
					Val:       val,
				}
				if decision.Qty > 0 {
					decisionArr = append(decisionArr, decision)
				}
			}
		}
	}

	for val, qty := range positions {
		currentProfit := 1 - ((val / prediction.CloseValue) * (1 - fee))
		if currentProfit > s.Config.ProfitCap {
			decision := trader.Decision{
				EventType: trader.PROFIT_SELL,
				Coin:      prediction.Coin,
				Qty:       s.SellSize(prediction, qty),
				Val:       val,
			}
			if decision.Qty > 0 {
				decisionArr = append(decisionArr, decision)
			}
		}
	}

	if len(decisionArr) == 0 {
		decisionArr = append(decisionArr, trader.Decision{
			EventType: trader.HOLD,
			Coin:      prediction.Coin,
			Qty:       0,
			Val:       0,
		})
	}

	return decisionArr
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
