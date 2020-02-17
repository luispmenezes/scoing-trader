package strategies

import (
	"scoing-trader/trader/model/market/model"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
)

type BasicStrategy struct {
	Config BasicConfig
}

func NewBasicStrategy(slice []float64) *BasicStrategy {
	basicStrategy := &BasicStrategy{Config: BasicConfig{}}
	basicStrategy.Config.FromSlice(slice)
	return basicStrategy
}

func (s *BasicStrategy) ComputeDecision(prediction predictor.Prediction, positions map[int64]float64,
	coinNetWorth int64, coinValue int64, totalNetWorth int64, balance int64, fee float64) map[trader.DecisionType]trader.Decision {

	decisionMap := make(map[trader.DecisionType]trader.Decision)

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
			Val:       model.FloatToInt(prediction.CloseValue),
		}
		if decision.Qty > 0 {
			decisionMap[trader.BUY] = decision
		}
	}

	if ((pred15 * s.Config.SellPred5Mod) + (pred60 * s.Config.SellPred10Mod) + (pred1440 * s.Config.SellPred100Mod)) < -2 {
		for val, qty := range positions {
			currentProfit := 1 - ((model.IntToFloat(val) / prediction.CloseValue) * (1 - fee))
			if currentProfit < s.Config.StopLoss {
				decision := trader.Decision{
					EventType: trader.SELL,
					Coin:      prediction.Coin,
					Qty:       s.SellSize(prediction, qty, coinValue),
					Val:       val,
				}
				if decision.Qty > 0 {
					if sellDecision, exists := decisionMap[trader.SELL]; exists {
						sellDecision.Qty += decision.Qty
					} else {
						decisionMap[trader.SELL] = decision
					}
				}
			}
		}
	}

	for val, qty := range positions {
		currentProfit := 1 - ((model.IntToFloat(val) / prediction.CloseValue) * (1 - fee))
		if currentProfit > s.Config.ProfitCap {
			decision := trader.Decision{
				EventType: trader.SELL,
				Coin:      prediction.Coin,
				Qty:       s.SellSize(prediction, qty, coinValue),
				Val:       val,
			}
			if decision.Qty > 0 {
				if sellDecision, exists := decisionMap[trader.SELL]; exists {
					sellDecision.Qty += decision.Qty
				} else {
					decisionMap[trader.SELL] = decision
				}
			}
		}
	}

	if len(decisionMap) == 0 {
		decisionMap[trader.HOLD] = trader.Decision{
			EventType: trader.HOLD,
			Coin:      prediction.Coin,
			Qty:       0,
			Val:       0,
		}
	}

	return decisionMap
}

func (s *BasicStrategy) BuySize(prediction predictor.Prediction, coinNetWorth int64, totalNetWorth int64,
	balance int64, fee float64) float64 {

	maxCoinNetWorth := model.IntFloatMul(totalNetWorth, 0.3)
	maxTransaction := model.IntFloatMul(totalNetWorth, 0.05)

	if maxCoinNetWorth-coinNetWorth >= 10 && balance >= model.FloatToInt(10*(1+fee)) {
		transaction := model.Max(10, model.Min(maxTransaction, model.IntFloatMul(maxCoinNetWorth-coinNetWorth, s.Config.BuyQtyMod)))
		transactionWFee := model.IntFloatMul(transaction, 1+fee)

		if transactionWFee < balance {
			return model.IntToFloat(transaction) / prediction.CloseValue
		} else {
			return model.IntToFloat(balance) / prediction.CloseValue * (1 + fee)
		}
	} else {
		return 0.0
	}
}

func (s *BasicStrategy) SellSize(prediction predictor.Prediction, positionQty float64, coinValue int64) float64 {
	proposedQty := positionQty * s.Config.SellQtyMod

	if model.IntFloatMul(coinValue, proposedQty) > 10 {
		return proposedQty
	} else {
		return 0.0
	}
}
