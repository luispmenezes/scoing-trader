package strategies

import (
	"math"
	"scoing-trader/trader/model/market/model"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
)

type BasicWithMemoryStrategy struct {
	Config             BasicWithMemoryConfig
	PriceHistory       []float64
	PredictionHistory5 []float64
	DecisionHistory    []trader.DecisionType
	HistoryLength      int
}

func NewBasicWithMemoryStrategy(slice []float64, historyLength int) *BasicWithMemoryStrategy {
	basicWithMemoryStrategy := &BasicWithMemoryStrategy{
		Config:             BasicWithMemoryConfig{},
		PriceHistory:       make([]float64, 0),
		PredictionHistory5: make([]float64, 0),
		DecisionHistory:    make([]trader.DecisionType, 0),
		HistoryLength:      historyLength,
	}
	basicWithMemoryStrategy.Config.FromSlice(slice)
	return basicWithMemoryStrategy
}

func (s *BasicWithMemoryStrategy) ComputeDecision(prediction predictor.Prediction, positions map[int64]float64,
	coinNetWorth int64, coinValue int64, totalNetWorth int64, balance int64, fee float64) map[trader.DecisionType]trader.Decision {

	decisionMap := make(map[trader.DecisionType]trader.Decision)

	pred5, pred10, pred100 := s.computePredictors(prediction.Pred5, prediction.Pred10, prediction.Pred100)

	var priceDelta, predDelta float64

	if len(s.PriceHistory) >= s.HistoryLength {
		priceDelta = s.historyGetPriceDelta()
		predDelta = s.historyGetPred5Delta()
	}

	if len(s.PriceHistory) < s.HistoryLength || (s.historyGetDecisionCount(trader.BUY) <
		math.Round(float64(s.HistoryLength)/2) && priceDelta != -1 && predDelta != -1) {

		if ((pred5 * s.Config.BuyPred5Mod) + (pred10 * s.Config.BuyPred10Mod) + (pred100 * s.Config.BuyPred100Mod)) > 2 {
			if buyQty := s.BuySize(prediction, coinNetWorth, totalNetWorth, balance, fee); buyQty > 0 {
				decisionMap[trader.BUY] = trader.Decision{
					EventType: trader.BUY,
					Coin:      prediction.Coin,
					Qty:       buyQty,
					Val:       model.FloatToInt(prediction.CloseValue),
				}
			}
		}
	}

	if len(s.PriceHistory) < s.HistoryLength || (s.historyGetDecisionCount(trader.SELL) <
		math.Round(float64(s.HistoryLength)/2) && priceDelta != 1 && predDelta != 1) {

		if ((pred5 * s.Config.SellPred5Mod) + (pred10 * s.Config.SellPred10Mod) + (pred100 * s.Config.SellPred100Mod)) < -2 {
			for val, qty := range positions {
				currentProfit := 1 - ((model.IntToFloat(val) / prediction.CloseValue) * (1 - fee))
				if currentProfit < s.Config.StopLoss {
					if sellQty := s.SellSize(prediction, qty, coinValue); sellQty > 0 {
						if sellDecision, exists := decisionMap[trader.SELL]; exists {
							sellDecision.Qty += sellQty
						} else {
							decisionMap[trader.SELL] = trader.Decision{
								EventType: trader.SELL,
								Coin:      prediction.Coin,
								Qty:       sellQty,
								Val:       val,
							}
						}
					}
				}
			}
		}
	}
	for val, qty := range positions {
		currentProfit := 1 - ((model.IntToFloat(val) / prediction.CloseValue) * (1 - fee))
		if currentProfit > s.Config.ProfitCap {
			if sellQty := s.SellSize(prediction, qty, coinValue); sellQty > 0 {
				if sellDecision, exists := decisionMap[trader.SELL]; exists {
					sellDecision.Qty += sellQty
				} else {
					decisionMap[trader.SELL] = trader.Decision{
						EventType: trader.SELL,
						Coin:      prediction.Coin,
						Qty:       sellQty,
						Val:       val,
					}
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

	decisionTypes := make([]trader.DecisionType, 0, len(decisionMap))
	for d := range decisionMap {
		decisionTypes = append(decisionTypes, d)
	}

	s.addToHistory(prediction.CloseValue, prediction.Pred5, decisionTypes)

	return decisionMap
}

func (s *BasicWithMemoryStrategy) BuySize(prediction predictor.Prediction, coinNetWorth int64, totalNetWorth int64,
	balance int64, fee float64) float64 {

	maxCoinNetWorth := model.IntFloatMul(totalNetWorth, 0.3)
	maxTransaction := model.IntFloatMul(totalNetWorth, 0.05)

	if maxCoinNetWorth-coinNetWorth >= 1000000000 && balance >= model.FloatToInt(10*(1+fee)) {
		transaction := model.Max(1000000000, model.Min(maxTransaction, model.IntFloatMul(maxCoinNetWorth-coinNetWorth, s.Config.BuyQtyMod)))
		transactionWFee := model.IntFloatMul(transaction, 1+fee)

		if transactionWFee < balance {
			return model.IntToFloat(transaction) / prediction.CloseValue
		} else {
			return model.IntToFloat(balance-1) / (prediction.CloseValue * (1 + fee))
		}
	} else {
		return 0.0
	}
}

func (s *BasicWithMemoryStrategy) SellSize(prediction predictor.Prediction, positionQty float64, coinValue int64) float64 {
	proposedQty := positionQty * s.Config.SellQtyMod

	if model.IntFloatMul(coinValue, proposedQty) > 1000000000 {
		return proposedQty
	} else {
		return 0.0
	}
}

func (s *BasicWithMemoryStrategy) computePredictors(predValue5 float64, predValue10 float64, predValue100 float64) (float64, float64, float64) {
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

func (s *BasicWithMemoryStrategy) addToHistory(price float64, pred5 float64, decisionTypes []trader.DecisionType) {
	s.PriceHistory = append([]float64{price}, s.PriceHistory...)
	s.PredictionHistory5 = append([]float64{pred5}, s.PredictionHistory5...)
	s.DecisionHistory = append(decisionTypes, s.DecisionHistory...)

	if len(s.PriceHistory) > s.HistoryLength {
		s.PriceHistory = s.PriceHistory[:s.HistoryLength]
		s.PredictionHistory5 = s.PredictionHistory5[:s.HistoryLength]
		s.DecisionHistory = s.DecisionHistory[:s.HistoryLength]
	}
}

func (s *BasicWithMemoryStrategy) historyGetDecisionCount(decisionType trader.DecisionType) float64 {
	total := 0.0

	for _, dec := range s.DecisionHistory {
		if dec == decisionType {
			total++
		}
	}

	return total
}

func (s *BasicWithMemoryStrategy) historyGetPriceDelta() float64 {
	absDelta := (s.PriceHistory[len(s.PriceHistory)-1] - s.PriceHistory[0]) / math.Abs(s.PriceHistory[0])

	if absDelta > 0.05 {
		return 1
	} else if absDelta < -0.05 {
		return -1
	} else {
		return 0
	}
}

func (s *BasicWithMemoryStrategy) historyGetPred5Delta() float64 {
	absDelta := (s.PredictionHistory5[len(s.PredictionHistory5)-1] - s.PredictionHistory5[0]) / math.Abs(s.PredictionHistory5[0])

	if absDelta > 0.05 {
		return 1
	} else if absDelta < -0.05 {
		return -1
	} else {
		return 0
	}
}
