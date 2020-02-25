package strategies

import (
	"github.com/shopspring/decimal"
	"math"
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

func (s *BasicWithMemoryStrategy) ComputeDecision(prediction predictor.Prediction, positions map[string]decimal.Decimal,
	coinNetWorth decimal.Decimal, coinValue decimal.Decimal, totalNetWorth decimal.Decimal, balance decimal.Decimal, fee decimal.Decimal) map[trader.DecisionType]trader.Decision {

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
			if buyQty := s.BuySize(prediction, coinNetWorth, totalNetWorth, balance, fee); buyQty.GreaterThan(decimal.Zero) {
				decisionMap[trader.BUY] = trader.Decision{
					EventType: trader.BUY,
					Coin:      prediction.Coin,
					Qty:       buyQty,
					Val:       decimal.NewFromFloat(prediction.CloseValue),
				}
			}
		}
	}

	if len(s.PriceHistory) < s.HistoryLength || (s.historyGetDecisionCount(trader.SELL) <
		math.Round(float64(s.HistoryLength)/2) && priceDelta != 1 && predDelta != 1) {

		if ((pred5 * s.Config.SellPred5Mod) + (pred10 * s.Config.SellPred10Mod) + (pred100 * s.Config.SellPred100Mod)) < -2 {
			for val, qty := range positions {
				decimalVal, _ := decimal.NewFromString(val)
				currentProfit := decimal.NewFromInt(1).Sub(decimalVal.Div(decimal.NewFromFloat(prediction.CloseValue)).Mul(decimal.NewFromInt(1).Sub(fee)))
				if (len(s.PriceHistory) < s.HistoryLength || (s.historyGetDecisionCount(trader.SELL) <
					math.Round(float64(s.HistoryLength)/2) && priceDelta != 1 && predDelta != 1) &&
					((pred5*s.Config.SellPred5Mod)+(pred10*s.Config.SellPred10Mod)+(pred100*s.Config.SellPred100Mod)) < -2 &&
					currentProfit.LessThan(decimal.NewFromFloat(s.Config.StopLoss))) ||
					currentProfit.GreaterThan(decimal.NewFromFloat(s.Config.ProfitCap)) {
					if sellQty := s.SellSize(prediction, qty, coinValue); sellQty.GreaterThan(decimal.Zero) {
						if sellDecision, exists := decisionMap[trader.SELL]; exists {
							sellDecision.Qty = sellDecision.Qty.Add(sellQty)
						} else {
							decisionMap[trader.SELL] = trader.Decision{
								EventType: trader.SELL,
								Coin:      prediction.Coin,
								Qty:       sellQty,
								Val:       decimalVal,
							}
						}
					}
				}
			}
		}
	}

	if len(decisionMap) == 0 {
		decisionMap[trader.HOLD] = trader.Decision{
			EventType: trader.HOLD,
			Coin:      prediction.Coin,
			Qty:       decimal.Zero,
			Val:       decimal.Zero,
		}
	}

	decisionTypes := make([]trader.DecisionType, 0, len(decisionMap))
	for d := range decisionMap {
		decisionTypes = append(decisionTypes, d)
	}

	s.addToHistory(prediction.CloseValue, prediction.Pred5, decisionTypes)

	return decisionMap
}

func (s *BasicWithMemoryStrategy) BuySize(prediction predictor.Prediction, coinNetWorth decimal.Decimal, totalNetWorth decimal.Decimal,
	balance decimal.Decimal, fee decimal.Decimal) decimal.Decimal {

	maxCoinNetWorth := totalNetWorth.Mul(decimal.NewFromFloat(0.3))
	maxTransaction := totalNetWorth.Mul(decimal.NewFromFloat(0.05))

	if maxCoinNetWorth.Sub(coinNetWorth).GreaterThanOrEqual(decimal.NewFromInt(10)) &&
		balance.GreaterThanOrEqual(decimal.NewFromInt(10).Mul(decimal.NewFromInt(1).Mul(fee))) {
		transaction := decimal.Max(decimal.NewFromInt(10), decimal.Min(maxTransaction, maxCoinNetWorth.Sub(coinNetWorth).Mul(decimal.NewFromFloat(s.Config.BuyQtyMod))))
		transactionWFee := transaction.Mul(decimal.NewFromInt(1).Add(fee))

		if transactionWFee.LessThan(balance) {
			return transaction.Div(decimal.NewFromFloat(prediction.CloseValue))
		} else {
			return balance.Sub(decimal.NewFromInt(1)).Div(decimal.NewFromFloat(prediction.CloseValue).Mul(decimal.NewFromInt(1).Add(fee)))
		}
	} else {
		return decimal.Zero
	}
}

func (s *BasicWithMemoryStrategy) SellSize(prediction predictor.Prediction, positionQty decimal.Decimal, coinValue decimal.Decimal) decimal.Decimal {
	proposedQty := positionQty.Mul(decimal.NewFromFloat(s.Config.SellQtyMod))

	if coinValue.Mul(proposedQty).GreaterThan(decimal.NewFromInt(10)) {
		return proposedQty
	} else {
		return decimal.Zero
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
