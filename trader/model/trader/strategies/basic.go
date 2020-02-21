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

	pred5, pred10, pred100 := s.computePredictors(prediction.Pred5, prediction.Pred10, prediction.Pred100)

	if ((pred5 * s.Config.BuyPred5Mod) + (pred10 * s.Config.BuyPred10Mod) + (pred100 * s.Config.BuyPred100Mod)) > 2 {
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

	if ((pred5 * s.Config.SellPred5Mod) + (pred10 * s.Config.SellPred10Mod) + (pred100 * s.Config.SellPred100Mod)) < -2 {
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
			return model.IntToFloat(balance) / (prediction.CloseValue * (1 + fee))
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

func (s *BasicStrategy) computePredictors(predValue5 float64, predValue10 float64, predValue100 float64) (float64, float64, float64) {
	pred5 := 0.0
	pred10 := 0.0
	pred100 := 0.0

	if predValue5 > 0.01 {
		pred5 = 1
	} else if predValue5 < -0.01 {
		pred5 = -1
	}

	if predValue10 > 0.01 {
		pred10 = 1
	} else if predValue10 < -0.01 {
		pred10 = -1
	}

	if predValue100 > 0.01 {
		pred100 = 1
	} else if predValue100 < -0.01 {
		pred100 = -1
	}

	return pred5, pred10, pred100
}
