package predictor

import "time"

type Predictor interface {
	Predict(coin string) Prediction
	SetNextPrediction(prediction Prediction)
}

type Prediction struct {
	Timestamp  time.Time `json:"open_time"`
	Coin       string    `json:"coin"`
	CloseValue float64   `json:"close_value"`
	Pred15     float64   `json:"pred_15"`
	Pred60     float64   `json:"pred_60"`
	Pred1440   float64   `json:"pred_1440"`
}
