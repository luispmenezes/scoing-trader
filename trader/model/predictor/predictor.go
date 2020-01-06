package predictor

import (
	"super-trader/trader/model"
	"time"
)

type Predictor interface {
	Predict(coin string, data model.ExchangeData) float64
}

type Prediction struct {
	Timestamp 	   time.Time
	Coin	 	   string
	PredictedValue float64
}