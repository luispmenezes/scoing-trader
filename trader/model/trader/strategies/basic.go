package strategies

import (
	"github.com/shopspring/decimal"
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

func (s *BasicStrategy) ComputeDecision(prediction predictor.Prediction, positions map[string]decimal.Decimal,
	coinNetWorth decimal.Decimal, coinValue decimal.Decimal, totalNetWorth decimal.Decimal, balance decimal.Decimal, fee decimal.Decimal) map[trader.DecisionType]trader.Decision {

	decisionMap := make(map[trader.DecisionType]trader.Decision)

	pred5, pred10, pred100 := s.computePredictors(prediction.Pred5, prediction.Pred10, prediction.Pred100)

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

	for val, qty := range positions {
		decimalVal, _ := decimal.NewFromString(val)
		currentProfit := decimal.NewFromInt(1).Sub(decimalVal.Div(decimal.NewFromFloat(prediction.CloseValue)).Mul(decimal.NewFromInt(1).Sub(fee)))
		if (((pred5*s.Config.SellPred5Mod)+(pred10*s.Config.SellPred10Mod)+(pred100*s.Config.SellPred100Mod)) < -2 &&
			currentProfit.LessThan(decimal.NewFromFloat(s.Config.StopLoss))) || currentProfit.GreaterThan(decimal.NewFromFloat(s.Config.ProfitCap)) {
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

	if len(decisionMap) == 0 {
		decisionMap[trader.HOLD] = trader.Decision{
			EventType: trader.HOLD,
			Coin:      prediction.Coin,
			Qty:       decimal.Zero,
			Val:       decimal.Zero,
		}
	}

	return decisionMap
}

func (s *BasicStrategy) BuySize(prediction predictor.Prediction, coinNetWorth decimal.Decimal, totalNetWorth decimal.Decimal,
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

func (s *BasicStrategy) SellSize(prediction predictor.Prediction, positionQty decimal.Decimal, coinValue decimal.Decimal) decimal.Decimal {
	proposedQty := positionQty.Mul(decimal.NewFromFloat(s.Config.SellQtyMod))

	if coinValue.Mul(proposedQty).GreaterThan(decimal.NewFromInt(10)) {
		return proposedQty
	} else {
		return decimal.Zero
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
