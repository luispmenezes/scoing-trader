package strategies

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
)

type BasicWithMemoryStrategy struct {
	Config             BasicWithMemoryConfig
	PriceHistory       map[string][]float64
	PredictionHistory5 map[string][]float64
	DecisionHistory    map[string][]trader.DecisionType
	HistoryLength      int
}

func NewBasicWithMemoryStrategy(slice []float64, historyLength int) *BasicWithMemoryStrategy {
	basicWithMemoryStrategy := &BasicWithMemoryStrategy{
		Config:             BasicWithMemoryConfig{},
		PriceHistory:       make(map[string][]float64, 0),
		PredictionHistory5: make(map[string][]float64, 0),
		DecisionHistory:    make(map[string][]trader.DecisionType, 0),
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

	if len(s.PriceHistory[prediction.Coin]) >= s.HistoryLength {
		priceDelta = s.historyGetPriceDelta(prediction.Coin)
		predDelta = s.historyGetPred5Delta(prediction.Coin)
	}

	debugText := "\n\tDecision History: "

	for _, decision := range s.DecisionHistory[prediction.Coin] {
		debugText += string(decision[0]) + " "
	}

	debugText += fmt.Sprintf("\n\tPriceDelta: %.2f Predicition Delta: %.2f", priceDelta, predDelta)

	if len(s.PriceHistory[prediction.Coin]) < s.HistoryLength || (s.historyGetDecisionCount(prediction.Coin, trader.BUY) <
		math.Round(float64(s.HistoryLength)/2) && priceDelta != -1 && predDelta != -1) {

		if ((pred5 * s.Config.BuyPred5Mod) + (pred10 * s.Config.BuyPred10Mod) + (pred100 * s.Config.BuyPred100Mod)) > 2 {
			if buyQty := s.BuySize(prediction, coinNetWorth, totalNetWorth, balance, fee); buyQty.GreaterThan(decimal.Zero) {
				decisionMap[trader.BUY] = trader.Decision{
					EventType: trader.BUY,
					Coin:      prediction.Coin,
					Qty:       buyQty,
					BuyConf:   1,
					DebugText: debugText,
				}
			}
		}
	}

	if len(s.PriceHistory[prediction.Coin]) < s.HistoryLength || (s.historyGetDecisionCount(prediction.Coin, trader.SELL) <
		math.Round(float64(s.HistoryLength)/2) && priceDelta != 1 && predDelta != 1) {

		for val, qty := range positions {
			decimalVal, _ := decimal.NewFromString(val)
			currentProfit := decimal.NewFromInt(1).Sub(decimalVal.Div(decimal.NewFromFloat(prediction.CloseValue)).Mul(decimal.NewFromInt(1).Sub(fee)))
			if (len(s.PriceHistory) < s.HistoryLength || (s.historyGetDecisionCount(prediction.Coin, trader.SELL) <
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
							SellConf:  1,
							DebugText: debugText,
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
			BuyConf:   (pred5 * s.Config.BuyPred5Mod) + (pred10 * s.Config.BuyPred10Mod) + (pred100 * s.Config.BuyPred100Mod),
			SellConf:  (pred5 * s.Config.SellPred5Mod) + (pred10 * s.Config.SellPred10Mod) + (pred100 * s.Config.SellPred100Mod),
			DebugText: debugText,
		}
	}

	decisionTypes := make([]trader.DecisionType, 0, len(decisionMap))
	for d := range decisionMap {
		decisionTypes = append(decisionTypes, d)
	}

	s.addToHistory(prediction.Coin, prediction.CloseValue, prediction.Pred5, decisionTypes)

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

	if predValue5 > s.Config.SegTh {
		pred5 = 1
	} else if predValue5 < -s.Config.SegTh {
		pred5 = -1
	}

	if predValue10 > s.Config.SegTh {
		pred10 = 1
	} else if predValue10 < -s.Config.SegTh {
		pred10 = -1
	}

	if predValue100 > s.Config.SegTh {
		pred100 = 1
	} else if predValue100 < -s.Config.SegTh {
		pred100 = -1
	}

	return pred5, pred10, pred100
}

func (s *BasicWithMemoryStrategy) addToHistory(coin string, price float64, pred5 float64, decisionTypes []trader.DecisionType) {
	s.PriceHistory[coin] = append([]float64{price}, s.PriceHistory[coin]...)
	s.PredictionHistory5[coin] = append([]float64{pred5}, s.PredictionHistory5[coin]...)
	s.DecisionHistory[coin] = append(decisionTypes, s.DecisionHistory[coin]...)

	if len(s.PriceHistory[coin]) > s.HistoryLength {
		s.PriceHistory[coin] = s.PriceHistory[coin][:s.HistoryLength]
		s.PredictionHistory5[coin] = s.PredictionHistory5[coin][:s.HistoryLength]
		s.DecisionHistory[coin] = s.DecisionHistory[coin][:s.HistoryLength]
	}
}

func (s *BasicWithMemoryStrategy) historyGetDecisionCount(coin string, decisionType trader.DecisionType) float64 {
	total := 0.0

	for _, dec := range s.DecisionHistory[coin] {
		if dec == decisionType {
			total++
		}
	}

	return total
}

func (s *BasicWithMemoryStrategy) historyGetPriceDelta(coin string) float64 {
	absDelta := (s.PriceHistory[coin][len(s.PriceHistory[coin])-1] - s.PriceHistory[coin][0]) / math.Abs(s.PriceHistory[coin][0])

	if absDelta > s.Config.HistSegTh {
		return 1
	} else if absDelta < -s.Config.HistSegTh {
		return -1
	} else {
		return 0
	}
}

func (s *BasicWithMemoryStrategy) historyGetPred5Delta(coin string) float64 {
	absDelta := (s.PredictionHistory5[coin][len(s.PredictionHistory5[coin])-1] - s.PredictionHistory5[coin][0]) / math.Abs(s.PredictionHistory5[coin][0])

	if absDelta > s.Config.HistSegTh {
		return 1
	} else if absDelta < -s.Config.HistSegTh {
		return -1
	} else {
		return 0
	}
}
