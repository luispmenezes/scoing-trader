package strategies

import (
	"math"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
)

type BasicStrategy struct {
	Config trader.TraderConfig
}

func NewBasicStrategy(config trader.TraderConfig) *BasicStrategy {
	return &BasicStrategy{Config: config}
}

func (s *BasicStrategy) ComputeDecision(prediction predictor.Prediction, positions map[float64]float64, coinNetWorth float64, totalNetWorth float64) trader.Decision {
	if ((s.Config.BuyPred15Mod * prediction.Pred15) + (s.Config.BuyPred60Mod * prediction.Pred60) +
		(s.Config.BuyPred1440Mod * prediction.Pred1440)) > 0 {
		if (coinNetWorth / totalNetWorth) < 0.3 {
			return trader.Decision{
				EventType: trader.BUY,
				Qty:       s.BuySize(prediction, coinNetWorth, totalNetWorth),
				Val:       prediction.OpenValue,
			}
		}
	}

	if ((s.Config.SellPred15Mod * prediction.Pred15) + (s.Config.SellPred60Mod * prediction.Pred60) +
		(s.Config.SellPred1440Mod * prediction.Pred1440)) > 0 {
		for val, qty := range positions {
			currentProfit := qty * prediction.OpenValue
			if currentProfit < s.Config.StopLoss || currentProfit > s.Config.ProfitCap {
				return trader.Decision{
					EventType: trader.SELL,
					Qty:       s.SellSize(prediction, qty),
					Val:       val,
				}
			}
		}
	}

	return trader.Decision{
		EventType: trader.HOLD,
		Qty:       0,
	}
}

func (s *BasicStrategy) BuySize(prediction predictor.Prediction, coinNetWorth float64, totalNetWorth float64) float64 {
	proposedSize := (coinNetWorth * s.Config.BuyNWQtyMod) + (prediction.Pred15 * s.Config.BuyQty15Mod) +
		(prediction.Pred60 * s.Config.BuyQty60Mod) + (prediction.Pred1440 * s.Config.BuyQty1440Mod)
	return math.Min(proposedSize, totalNetWorth*0.1)
}

func (s *BasicStrategy) SellSize(prediction predictor.Prediction, positionQty float64) float64 {
	proposedSize := (positionQty * s.Config.SellPosQtyMod) + (prediction.Pred15 * s.Config.SellQty15Mod) +
		(prediction.Pred60 * s.Config.SellQty60Mod) + (prediction.Pred1440 * s.Config.SellQty1440Mod)
	return math.Min(proposedSize, positionQty)
}
