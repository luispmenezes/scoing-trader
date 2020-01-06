package trader

import (
	"log"
	"super-trader/trader/model"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/wallet"
)

type Trader struct {
	Config    TraderConfig
	Wallet    wallet.Wallet
	Predictor predictor.Predictor
	Logging   bool
}

func NewTrader(config TraderConfig, wallet wallet.Wallet, predictor predictor.Predictor, logging bool) *Trader {
	return &Trader{
		Config:    config,
		Wallet:    wallet,
		Predictor: predictor,
		Logging:   logging,
	}
}

func (t *Trader) ProcessData(data model.ExchangeData, coin string) {
	var currentValue = data.CloseValue
	var predictionDelta = t.Predictor.Predict(coin,data)
	t.Wallet.UpdateCoinValue(coin, currentValue, data.OpenTime)

	if t.Logging {
		log.Printf("%v NW: %f B: %f Prediction: %f CurrentValue: %f (%f%%)", data.OpenTime, t.Wallet.NetWorth(),
			t.Wallet.GetBalance(), currentValue*predictionDelta, currentValue, predictionDelta)
	}
	if len(t.Wallet.GetPositions(coin)) == 0 && predictionDelta >= t.Config.BuyThreshold {
		t.Wallet.Buy(coin, t.getBuySizing()/currentValue)

		if t.Logging {
			log.Printf("BUY: %f at %f$", t.getBuySizing(), currentValue)
		}
	} else if predictionDelta >= t.Config.IncreaseThreshold {
		t.Wallet.Buy(coin, t.getBuySizing()/currentValue)

		if t.Logging {
			log.Printf("BUY: %f at %f$", t.getBuySizing(), currentValue)
		}
	} else if predictionDelta < t.Config.SellThreshold {
		for posValue, posSize := range t.Wallet.GetPositions(coin) {
			if (currentValue/posValue)-1 < t.Config.MaxLoss {
				if t.Logging {
					log.Printf("DUMP SELL: %f %s at %f$ ( Bought @ %f$ ) -- Profit %f", posSize, coin,
						currentValue, posValue, (currentValue-posValue)*posSize)
				}
				t.Wallet.Sell(coin, posValue, posSize)
			}
		}
	} else {
		for posValue, posSize := range t.Wallet.GetPositions(coin) {
			if (currentValue/posValue)-1 > t.Config.MinProfit {
				if t.Logging {
					log.Printf("PROFIT SELL: %f %s at %f$ ( Bought @ %f$ ) -- Profit %f", posSize, coin,
						currentValue, posValue, (currentValue-posValue)*posSize)
				}
				t.Wallet.Sell(coin, posValue, posSize)
			}
		}
	}
}

func (t *Trader) getBuySizing() float64 {
	return t.Config.PositionSizing * t.Wallet.GetBalance()
}

func (t *Trader) getIncreaseBuySizing() float64 {
	return t.Config.IncreaseSizing * t.Wallet.GetBalance()
}